package handlers

import (
	"encoding/json"
	"net/http"
)

func JWKSHandler(w http.ResponseWriter, r *http.Request) {
	// Serve public key as JWK Set
	jwks := map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"kty": "RSA",
				"kid": "your-key-id",
				// Add other public key information
			},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)
}
