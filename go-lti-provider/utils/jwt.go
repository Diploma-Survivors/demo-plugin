package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// Moodle local endpoints
const (
	moodleJWKSURL = "http://localhost:8888/mod/lti/certs.php"
)

// OAuth2 Token Response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// GetAccessToken retrieves OAuth2 access token for AGS
func GetAccessToken() (string, error) {
	tokenURL := "http://localhost:8888/mod/lti/token.php"
	clientID := "wAWXk7ifY0o9tCU"
	clientSecret := "your-client-secret"
	scope := "https://purl.imsglobal.org/spec/lti-ags/scope/score"

	data := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s&scope=%s",
		clientID, clientSecret, scope)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

// VerifyJWT verifies JWT token từ Moodle sử dụng JWKS
func VerifyJWT(tokenString string) (map[string]interface{}, error) {
	// Fetch JWKS từ Moodle
	keySet, err := jwk.Fetch(context.Background(), moodleJWKSURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS from %s: %w", moodleJWKSURL, err)
	}

	// Parse và verify JWT token
	token, err := jwt.Parse([]byte(tokenString), jwt.WithKeySet(keySet))
	if err != nil {
		return nil, fmt.Errorf("failed to parse and verify JWT: %w", err)
	}

	// Convert token claims to map
	claims := make(map[string]interface{})
	for key, value := range token.PrivateClaims() {
		claims[key] = value
	}

	// Add standard claims
	if iss := token.Issuer(); iss != "" {
		claims["iss"] = iss
	}
	if sub := token.Subject(); sub != "" {
		claims["sub"] = sub
	}
	if aud := token.Audience(); len(aud) > 0 {
		claims["aud"] = aud
	}
	if exp := token.Expiration(); !exp.IsZero() {
		claims["exp"] = exp.Unix()
	}
	if iat := token.IssuedAt(); !iat.IsZero() {
		claims["iat"] = iat.Unix()
	}
	if nbf := token.NotBefore(); !nbf.IsZero() {
		claims["nbf"] = nbf.Unix()
	}
	if jti := token.JwtID(); jti != "" {
		claims["jti"] = jti
	}

	return claims, nil
}

// VerifyJWTWithKey verifies JWT với specific key (for testing)
func VerifyJWTWithKey(tokenString string, key interface{}) (map[string]interface{}, error) {
	token, err := jwt.Parse([]byte(tokenString), jwt.WithKey(jwa.RS256, key))
	if err != nil {
		return nil, fmt.Errorf("failed to parse and verify JWT with key: %w", err)
	}

	// Convert token claims to map
	claims := make(map[string]interface{})
	for key, value := range token.PrivateClaims() {
		claims[key] = value
	}

	// Add standard claims
	if iss := token.Issuer(); iss != "" {
		claims["iss"] = iss
	}
	if sub := token.Subject(); sub != "" {
		claims["sub"] = sub
	}
	if aud := token.Audience(); len(aud) > 0 {
		claims["aud"] = aud
	}
	if exp := token.Expiration(); !exp.IsZero() {
		claims["exp"] = exp.Unix()
	}
	if iat := token.IssuedAt(); !iat.IsZero() {
		claims["iat"] = iat.Unix()
	}

	return claims, nil
}
