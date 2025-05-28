package handlers

import (
	"fmt"
	"go-lti-provider/utils"
	"log"
	"net/http"
	"net/url"
)

// LTILaunchRedirectHandler handles LTI 1.3 launch and redirects to frontend
func LTILaunchRedirectHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("🚀 LTI 1.3 Launch received - Redirecting to frontend")

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

	// Verify JWT and extract claims
	claims, err := utils.VerifyJWT(idToken)
	if err != nil {
		log.Printf("❌ JWT verification failed: %v", err)
		http.Error(w, "Invalid JWT token", http.StatusUnauthorized)
		return
	}

	// Extract required information
	userID := claims["sub"].(string)
	context := claims["https://purl.imsglobal.org/spec/lti/claim/context"].(map[string]interface{})
	contextID := context["id"].(string)

	// Extract AGS endpoint if available
	lineitem := ""
	if endpoint, ok := claims["https://purl.imsglobal.org/spec/lti-ags/claim/endpoint"].(map[string]interface{}); ok {
		if li, ok := endpoint["lineitem"].(string); ok {
			lineitem = li
		}
	}

	// Build frontend redirect URL with parameters
	feURL := fmt.Sprintf("http://localhost:3000?id_token=%s&user=%s&context=%s",
		url.QueryEscape(idToken),
		url.QueryEscape(userID),
		url.QueryEscape(contextID))

	// Add lineitem if available
	if lineitem != "" {
		feURL += fmt.Sprintf("&lineitem=%s", url.QueryEscape(lineitem))
	}

	log.Printf("✅ Redirecting to frontend: %s", feURL)
	http.Redirect(w, r, feURL, http.StatusSeeOther)
}
