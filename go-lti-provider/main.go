package main

import (
	"log"
	"net/http"
	"os"

	"go-lti-provider/config"
	"go-lti-provider/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set default port if not provided
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware(cfg.AllowedOrigins))

	// LTI routes
	r.Route("/lti", func(r chi.Router) {
		r.Post("/login", handlers.LoginHandler)
		r.Post("/launch", handlers.LTILaunchRedirectHandler)
		r.Post("/grade", handlers.GradeHandler)
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/execute", handlers.ExecuteHandler)
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// JWKS endpoint
	r.Get("/.well-known/jwks.json", handlers.JWKSHandler)

	// Start server
	log.Printf("🚀 LTI Provider: http://localhost:%s", port)
	log.Printf("⚡ Judge0: %s", cfg.GetJudge0SubmissionURL())

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// corsMiddleware adds CORS headers
func corsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// // formatJSON formats map as JSON string (simple implementation)
// func formatJSON(data map[string]interface{}) string {
// 	var parts []string
// 	parts = append(parts, "{")

// 	i := 0
// 	for key, value := range data {
// 		if i > 0 {
// 			parts = append(parts, ",")
// 		}
// 		parts = append(parts, `"`+key+`":"`+toString(value)+`"`)
// 		i++
// 	}

// 	parts = append(parts, "}")
// 	return strings.Join(parts, "")
// }

// func toString(v interface{}) string {
// 	switch val := v.(type) {
// 	case string:
// 		return val
// 	case int:
// 		return string(rune(val))
// 	default:
// 		return ""
// 	}
// }
