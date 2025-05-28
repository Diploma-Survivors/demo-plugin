package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"go-lti-provider/utils"
)

// LTI 1.3 Launch Claims
type LTILaunchClaims struct {
	Issuer             string                 `json:"iss"`
	Subject            string                 `json:"sub"`
	Audience           []string               `json:"aud"`
	ExpirationTime     int64                  `json:"exp"`
	IssuedAt           int64                  `json:"iat"`
	Nonce              string                 `json:"nonce"`
	MessageType        string                 `json:"https://purl.imsglobal.org/spec/lti/claim/message_type"`
	Version            string                 `json:"https://purl.imsglobal.org/spec/lti/claim/version"`
	DeploymentID       string                 `json:"https://purl.imsglobal.org/spec/lti/claim/deployment_id"`
	TargetLinkURI      string                 `json:"https://purl.imsglobal.org/spec/lti/claim/target_link_uri"`
	ResourceLink       ResourceLink           `json:"https://purl.imsglobal.org/spec/lti/claim/resource_link"`
	LaunchPresentation LaunchPresentation     `json:"https://purl.imsglobal.org/spec/lti/claim/launch_presentation"`
	Custom             map[string]interface{} `json:"https://purl.imsglobal.org/spec/lti/claim/custom,omitempty"`
	Context            Context                `json:"https://purl.imsglobal.org/spec/lti/claim/context"`
	ToolPlatform       ToolPlatform           `json:"https://purl.imsglobal.org/spec/lti/claim/tool_platform"`
	Roles              []string               `json:"https://purl.imsglobal.org/spec/lti/claim/roles"`
	// AGS Claims
	EndpointClaim *EndpointClaim `json:"https://purl.imsglobal.org/spec/lti-ags/claim/endpoint,omitempty"`
}

type ResourceLink struct {
	ID          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type LaunchPresentation struct {
	DocumentTarget string `json:"document_target,omitempty"`
	Height         int    `json:"height,omitempty"`
	Width          int    `json:"width,omitempty"`
	ReturnURL      string `json:"return_url,omitempty"`
}

type Context struct {
	ID    string   `json:"id"`
	Label string   `json:"label,omitempty"`
	Title string   `json:"title,omitempty"`
	Type  []string `json:"type,omitempty"`
}

type ToolPlatform struct {
	Name              string `json:"name,omitempty"`
	ContactEmail      string `json:"contact_email,omitempty"`
	Description       string `json:"description,omitempty"`
	URL               string `json:"url,omitempty"`
	ProductFamilyCode string `json:"product_family_code,omitempty"`
	Version           string `json:"version,omitempty"`
	GUID              string `json:"guid,omitempty"`
}

type EndpointClaim struct {
	Scope    []string `json:"scope,omitempty"`
	LineItem string   `json:"lineitem,omitempty"`
}

// Judge0 Integration
type Submission struct {
	SourceCode string `json:"source_code"`
	LanguageID int    `json:"language_id"`
	Stdin      string `json:"stdin,omitempty"`
}

type Judge0Response struct {
	Token         string  `json:"token,omitempty"`
	Status        Status  `json:"status,omitempty"`
	Stdout        *string `json:"stdout"`
	Stderr        *string `json:"stderr"`
	CompileOutput *string `json:"compile_output"`
	Message       *string `json:"message"`
	ExitCode      *int    `json:"exit_code"`
	ExitSignal    *int    `json:"exit_signal"`
	Time          *string `json:"time"`
	Memory        *int    `json:"memory"`
}

type Status struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

const (
	judge0URL = "http://localhost:2358/submissions"
)

// LaunchHandler xử lý LTI Launch request với JWT id_token
func LaunchHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("🚀 LTI 1.3 Launch received")

	// Parse form data
	if err := r.ParseForm(); err != nil {
		log.Printf("❌ Error parsing form: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get id_token from form
	idToken := r.FormValue("id_token")
	if idToken == "" {
		log.Println("❌ Missing id_token")
		http.Error(w, "Missing id_token", http.StatusBadRequest)
		return
	}

	// Verify JWT và extract claims
	claims, err := utils.VerifyJWT(idToken)
	if err != nil {
		log.Printf("❌ JWT verification failed: %v", err)
		http.Error(w, "Invalid JWT token", http.StatusUnauthorized)
		return
	}

	// Parse LTI claims
	var ltiClaims LTILaunchClaims
	claimsBytes, _ := json.Marshal(claims)
	if err := json.Unmarshal(claimsBytes, &ltiClaims); err != nil {
		log.Printf("❌ Error parsing LTI claims: %v", err)
		http.Error(w, "Invalid LTI claims", http.StatusBadRequest)
		return
	}

	log.Printf("✅ LTI Launch validated - User: %s, Resource: %s",
		ltiClaims.Subject, ltiClaims.ResourceLink.Title)

	// Extract custom parameters
	code, hasCode := ltiClaims.Custom["code"].(string)
	if !hasCode || code == "" {
		log.Println("⚠️ No custom code parameter found")
		renderSuccessPage(w, ltiClaims, "", "No code provided")
		return
	}

	// Get language from custom params, default to Go
	language, _ := ltiClaims.Custom["language"].(string)
	if language == "" {
		language = "go"
	}

	languageID := getLanguageID(language)

	// Submit to Judge0
	result, err := submitToJudge0(code, languageID)
	if err != nil {
		log.Printf("❌ Judge0 error: %v", err)
		renderSuccessPage(w, ltiClaims, code, fmt.Sprintf("Execution error: %v", err))
		return
	}

	// Render success page với kết quả
	renderSuccessPage(w, ltiClaims, code, formatJudge0Result(result))
}

func submitToJudge0(code string, languageID int) (*Judge0Response, error) {
	submission := Submission{
		SourceCode: code,
		LanguageID: languageID,
	}

	jsonData, err := json.Marshal(submission)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal submission: %w", err)
	}

	// Submit với wait=true để lấy kết quả ngay
	resp, err := http.Post(judge0URL+"?wait=true", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to POST to Judge0: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Judge0 returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result Judge0Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Judge0 response: %w", err)
	}

	return &result, nil
}

func formatJudge0Result(result *Judge0Response) string {
	if result == nil {
		return "No result"
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Status: %s\n", result.Status.Description))

	if result.Stdout != nil && *result.Stdout != "" {
		output.WriteString(fmt.Sprintf("Output:\n%s\n", *result.Stdout))
	}

	if result.Stderr != nil && *result.Stderr != "" {
		output.WriteString(fmt.Sprintf("Error:\n%s\n", *result.Stderr))
	}

	if result.CompileOutput != nil && *result.CompileOutput != "" {
		output.WriteString(fmt.Sprintf("Compile Output:\n%s\n", *result.CompileOutput))
	}

	if result.Time != nil {
		output.WriteString(fmt.Sprintf("Execution Time: %s\n", *result.Time))
	}

	if result.Memory != nil {
		output.WriteString(fmt.Sprintf("Memory Used: %d KB\n", *result.Memory))
	}

	return output.String()
}

func getLanguageID(language string) int {
	languageMap := map[string]int{
		"go":         75,
		"python":     71,
		"java":       62,
		"javascript": 63,
		"cpp":        54,
		"c":          50,
		"php":        68,
		"ruby":       72,
		"rust":       73,
		"swift":      83,
	}

	if id, exists := languageMap[language]; exists {
		return id
	}
	return 75 // Default to Go
}

func renderSuccessPage(w http.ResponseWriter, claims LTILaunchClaims, code, result string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>LTI 1.3 Launch Success</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { color: #2e7d32; border-bottom: 2px solid #4caf50; padding-bottom: 10px; margin-bottom: 20px; }
        .info-section { margin: 20px 0; padding: 15px; background: #f8f9fa; border-left: 4px solid #007bff; }
        .code-section { margin: 20px 0; }
        .code { background: #f4f4f4; padding: 15px; border-radius: 4px; font-family: monospace; white-space: pre-wrap; }
        .result { background: #e8f5e8; padding: 15px; border-radius: 4px; font-family: monospace; white-space: pre-wrap; }
        .return-btn { display: inline-block; background: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 4px; margin-top: 20px; }
        .return-btn:hover { background: #0056b3; }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="header">🚀 LTI 1.3 Launch Successful!</h1>
        
        <div class="info-section">
            <h3>📋 Launch Information</h3>
            <p><strong>User:</strong> %s</p>
            <p><strong>Resource:</strong> %s</p>
            <p><strong>Context:</strong> %s</p>
            <p><strong>Platform:</strong> %s</p>
        </div>

        %s

        %s

        %s
    </div>
</body>
</html>`,
		claims.Subject,
		claims.ResourceLink.Title,
		claims.Context.Title,
		claims.ToolPlatform.Name,
		renderCodeSection(code),
		renderResultSection(result),
		renderReturnLink(claims.LaunchPresentation.ReturnURL),
	)

	w.Write([]byte(html))
}

func renderCodeSection(code string) string {
	if code == "" {
		return ""
	}
	return fmt.Sprintf(`
        <div class="code-section">
            <h3>💻 Submitted Code</h3>
            <div class="code">%s</div>
        </div>`, code)
}

func renderResultSection(result string) string {
	if result == "" {
		return ""
	}
	return fmt.Sprintf(`
        <div class="code-section">
            <h3>⚡ Execution Result</h3>
            <div class="result">%s</div>
        </div>`, result)
}

func renderReturnLink(returnURL string) string {
	if returnURL == "" {
		return ""
	}
	return fmt.Sprintf(`<a href="%s" class="return-btn">← Return to Moodle</a>`, returnURL)
}
