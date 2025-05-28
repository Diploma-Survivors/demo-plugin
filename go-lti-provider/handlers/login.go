package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// LTI 1.3 OIDC Login Parameters
type OIDCLoginRequest struct {
	IssuerID        string `json:"iss"`
	LoginHint       string `json:"login_hint"`
	TargetLinkURI   string `json:"target_link_uri"`
	LTIMessageHint  string `json:"lti_message_hint,omitempty"`
	ClientID        string `json:"client_id,omitempty"`
	LTIDeploymentID string `json:"lti_deployment_id,omitempty"`
}

// LoginHandler xử lý OIDC Initiate Login request từ Moodle
// Đây là bước đầu tiên trong LTI 1.3 flow
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("🔐 LTI 1.3 OIDC Initiate Login received")

	// Parse form data
	if err := r.ParseForm(); err != nil {
		log.Printf("❌ Error parsing form: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Extract OIDC parameters
	loginReq := OIDCLoginRequest{
		IssuerID:        r.FormValue("iss"),
		LoginHint:       r.FormValue("login_hint"),
		TargetLinkURI:   r.FormValue("target_link_uri"),
		LTIMessageHint:  r.FormValue("lti_message_hint"),
		ClientID:        r.FormValue("client_id"),
		LTIDeploymentID: r.FormValue("lti_deployment_id"),
	}

	// Validate required parameters
	if loginReq.IssuerID == "" || loginReq.LoginHint == "" || loginReq.TargetLinkURI == "" {
		log.Printf("❌ Missing required OIDC parameters: iss=%s, login_hint=%s, target_link_uri=%s",
			loginReq.IssuerID, loginReq.LoginHint, loginReq.TargetLinkURI)
		http.Error(w, "Missing required OIDC parameters", http.StatusBadRequest)
		return
	}

	log.Printf("✅ OIDC Login - Issuer: %s, LoginHint: %s", loginReq.IssuerID, loginReq.LoginHint)

	// Build authorization URL để redirect về Moodle
	authURL, err := buildAuthorizationURL(loginReq)
	if err != nil {
		log.Printf("❌ Error building authorization URL: %v", err)
		http.Error(w, "Failed to build authorization URL", http.StatusInternalServerError)
		return
	}

	log.Printf("🔄 Redirecting to Moodle authorization: %s", authURL)

	// Redirect về Moodle với authorization request
	http.Redirect(w, r, authURL, http.StatusFound)
}

func buildAuthorizationURL(loginReq OIDCLoginRequest) (string, error) {
	// Moodle local auth endpoint
	authEndpoint := "http://localhost:8888/mod/lti/auth.php"

	params := url.Values{
		"response_type":    {"id_token"},
		"response_mode":    {"form_post"},
		"scope":            {"openid"},
		"client_id":        {"wAWXk7ifY0o9tCU"},                  // LTI Tool client ID từ Moodle
		"redirect_uri":     {"http://localhost:8080/lti/launch"}, // Launch URL của tool
		"login_hint":       {loginReq.LoginHint},
		"state":            {generateState()}, // Random state for security
		"nonce":            {generateNonce()}, // Random nonce for security
		"prompt":           {"none"},
		"lti_message_hint": {loginReq.LTIMessageHint},
	}

	return fmt.Sprintf("%s?%s", authEndpoint, params.Encode()), nil
}

// Utility functions - trong production nên dùng crypto/rand
func generateState() string {
	return "random-state-123" // Demo only
}

func generateNonce() string {
	return "random-nonce-456" // Demo only
}
