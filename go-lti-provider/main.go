package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const (
	consumerKey    = "DSA_2025"
	consumerSecret = "DSA_Secret_2025"
	judge0URL      = "http://localhost:2358/submissions" // Ensure Judge0 is accessible from this LTI provider
)

type Submission struct {
	SourceCode string `json:"source_code"`
	LanguageID int    `json:"language_id"`
}

type LTIRequest struct {
	Params      url.Values
	Headers     map[string]string // Note: This effectively holds all form parameters
	Request     *http.Request
	ConsumerKey string
}

// rfc3986Encode percent-encodes a string according to RFC3986,
// ensuring spaces are encoded as %20.
func rfc3986Encode(s string) string {
	return strings.ReplaceAll(url.QueryEscape(s), "+", "%20")
}

func main() {
	http.HandleFunc("/lti", ltiHandler) // Ensure Moodle's Tool URL path matches this exactly (case-sensitive)
	fmt.Println("LTI Provider started on :8081")
	http.ListenAndServe(":8081", nil)
}

func NewLTIRequest(r *http.Request) (*LTIRequest, error) {
	// For LTI, all parameters are typically in the POST body as x-www-form-urlencoded
	err := r.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("failed to parse form: %w", err)
	}

	// Extract LTI parameters into Headers map for convenience
	// r.Form will contain all parsed parameters (from query and body)
	ltiFormParams := make(map[string]string)
	for key, values := range r.Form {
		if len(values) > 0 {
			ltiFormParams[key] = values[0] // Taking the first value if multiple are present
		}
	}

	consumerKeyFromParam := r.Form.Get("oauth_consumer_key")
	if consumerKeyFromParam == "" {
		return nil, fmt.Errorf("missing oauth_consumer_key")
	}

	return &LTIRequest{
		Params:      r.Form,        // Keep original url.Values for multi-value params if needed by signature
		Headers:     ltiFormParams, // Simplified map for easy access (like custom_param)
		Request:     r,
		ConsumerKey: consumerKeyFromParam,
	}, nil
}

func (lr *LTIRequest) ValidateRequest(secret string) (bool, error) {
	if lr.ConsumerKey != consumerKey { // Compare with the global constant
		return false, fmt.Errorf("invalid consumer key: got %s, expected %s", lr.ConsumerKey, consumerKey)
	}

	signature := lr.Params.Get("oauth_signature")
	if signature == "" {
		return false, fmt.Errorf("missing oauth_signature")
	}

	paramsCopy := make(url.Values)
	for k, v := range lr.Params {
		if k != "oauth_signature" {
			paramsCopy[k] = v
			// fmt.Printf("Raw param %s: %s\n", k, v[0]) // Logging all values can be noisy
		}
	}

	method := lr.Request.Method

	// Construct the base URL from the request
	// Scheme and host should be lowercase, path is case-sensitive.
	scheme := lr.Request.URL.Scheme
	if scheme == "" {
		if lr.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	// Use r.Host as it includes the port if non-standard. Moodle will use this.
	// The LTI spec requires the host in the base string URI to be lowercase.
	host := strings.ToLower(lr.Request.Host)
	path := lr.Request.URL.Path // Path is case-sensitive

	reconstructedBaseURL := fmt.Sprintf("%s://%s%s", scheme, host, path)
	fmt.Printf("Reconstructed Base URL for sig calc: %s\n", reconstructedBaseURL)

	paramsList := make([]string, 0)
	for k, values := range paramsCopy {
		// OAuth parameters can have multiple values with the same key,
		// though it's less common in LTI launches. The spec implies
		// each key-value pair is distinct if values differ.
		// For simplicity matching most LTI, if a key has multiple values,
		// they should be sorted and added if that's how Moodle does it.
		// Moodle typically sends single values for LTI launch params.
		// The current loop correctly handles multiple values if they exist in lr.Params.
		for _, v := range values {
			encodedKey := rfc3986Encode(k)   // Use corrected encoder
			encodedValue := rfc3986Encode(v) // Use corrected encoder
			paramsList = append(paramsList, fmt.Sprintf("%s=%s", encodedKey, encodedValue))
		}
	}
	sort.Strings(paramsList) // Sorts "key=value" strings

	normalizedParams := strings.Join(paramsList, "&")

	baseString := fmt.Sprintf("%s&%s&%s",
		method,                              // HTTP method (e.g., "POST") - should not be escaped if simple
		rfc3986Encode(reconstructedBaseURL), // Base URL, RFC3986 encoded
		rfc3986Encode(normalizedParams))     // Normalized params, RFC3986 encoded

	// The signing key is consumer_secret + "&" + token_secret (empty for LTI 1.1)
	// Both parts should be RFC3986 encoded.
	signingKey := fmt.Sprintf("%s&", rfc3986Encode(secret))

	h := hmac.New(sha1.New, []byte(signingKey))
	h.Write([]byte(baseString))
	calculatedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	fmt.Printf("Method: %s\n", method)
	fmt.Printf("Base URL (used in sig): %s\n", reconstructedBaseURL)
	fmt.Printf("Normalized Parameters (before final encoding for base string): %s\n", normalizedParams)
	fmt.Printf("Signature Base String: %s\n", baseString)
	fmt.Printf("Signing Key: %s\n", signingKey)
	fmt.Printf("Calculated signature: %s\n", calculatedSignature)
	fmt.Printf("Received signature: %s\n", signature)

	return calculatedSignature == signature, nil
}

func (lr *LTIRequest) CreateReturnURL(errorMessage string) (*url.URL, error) {
	returnURLStr := lr.Params.Get("launch_presentation_return_url")
	if returnURLStr == "" {
		// No return URL, cannot redirect with error message
		return nil, fmt.Errorf("no launch_presentation_return_url provided")
	}

	parsedURL, err := url.Parse(returnURLStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse return URL: %w", err)
	}

	q := parsedURL.Query()
	if errorMessage != "" {
		q.Add("lti_msg", "LTI Launch Error") // General message
		q.Add("lti_errormsg", errorMessage)  // Specific error
	}
	// You might want to add lti_errorlog for more detailed internal logging too
	parsedURL.RawQuery = q.Encode()

	return parsedURL, nil
}

func ltiHandler(w http.ResponseWriter, r *http.Request) {
	ltiRequest, err := NewLTIRequest(r)
	if err != nil {
		// This error is before we can reliably get a launch_presentation_return_url
		http.Error(w, "Invalid LTI request structure: "+err.Error(), http.StatusBadRequest)
		return
	}

	valid, err := ltiRequest.ValidateRequest(consumerSecret)
	if err != nil { // This err is from ValidateRequest if something went wrong building its components
		fmt.Printf("Error during validation logic: %v\n", err) // Log internal error
		returnUrl, _ := ltiRequest.CreateReturnURL("Error during signature validation process.")
		if returnUrl != nil {
			http.Redirect(w, r, returnUrl.String(), http.StatusSeeOther) // Use 303 See Other for POST redirect
		} else {
			http.Error(w, "Error during signature validation process.", http.StatusInternalServerError)
		}
		return
	}

	if !valid {
		fmt.Println("LTI signature validation failed.")
		returnUrl, _ := ltiRequest.CreateReturnURL("LTI signature verification failed.")
		if returnUrl != nil {
			http.Redirect(w, r, returnUrl.String(), http.StatusSeeOther)
		} else {
			http.Error(w, "Invalid LTI signature.", http.StatusUnauthorized)
		}
		return
	}

	// Signature is valid, proceed with LTI launch
	fmt.Println("LTI signature validation successful.")

	// Correctly extract the custom parameter (Moodle prefixes with "custom_")
	code, ok := ltiRequest.Headers["custom_custom_code"] // Moodle sends "custom_YOUR_PARAM_NAME"
	if !ok {
		errorMsg := "Required custom parameter 'custom_custom_code' not found."
		fmt.Println(errorMsg)
		returnUrl, _ := ltiRequest.CreateReturnURL(errorMsg)
		if returnUrl != nil {
			http.Redirect(w, r, returnUrl.String(), http.StatusSeeOther)
		} else {
			http.Error(w, errorMsg, http.StatusBadRequest)
		}
		return
	}

	// Submit code to Judge0
	result, err := submitToJudge0(code, 63) // Assuming 63 is Go. Check Judge0 docs for language IDs.
	if err != nil {
		errorMsg := fmt.Sprintf("Error executing code via Judge0: %v", err)
		fmt.Println(errorMsg)
		// Decide if this error should go back to Moodle UI or just be logged
		returnUrl, _ := ltiRequest.CreateReturnURL(errorMsg)
		if returnUrl != nil {
			http.Redirect(w, r, returnUrl.String(), http.StatusSeeOther)
		} else {
			http.Error(w, errorMsg, http.StatusInternalServerError)
		}
		return
	}

	// Successfully processed
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>LTI Launch Successful!</h1><p>Code was: %s</p><p>Judge0 Result: <pre>%s</pre></p>", code, result)
	fmt.Fprintf(w, "<p><a href=\"%s\">Return to Moodle</a></p>", ltiRequest.Params.Get("launch_presentation_return_url"))

}

func submitToJudge0(code string, languageID int) (string, error) {
	submission := Submission{SourceCode: code, LanguageID: languageID}
	jsonData, err := json.Marshal(submission)
	if err != nil {
		return "", fmt.Errorf("failed to marshal submission JSON: %w", err)
	}

	resp, err := http.Post(judge0URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to POST to Judge0: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 { // Check for non-2xx status codes from Judge0
		// Attempt to read body for more info, but it might be empty or non-JSON
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Judge0 returned error status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode Judge0 response JSON: %w", err)
	}

	// For better display, marshal the result map back to a formatted JSON string
	prettyResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", result), nil // Fallback to default map string representation
	}
	return string(prettyResult), nil
}
