package handlers

import (
	"encoding/json"
	"fmt"
	"go-lti-provider/config"
	"go-lti-provider/models"
	"go-lti-provider/services"
	"log"
	"net/http"
	"time"
)

// ExecuteRequest represents the request body for code execution
type ExecuteRequest struct {
	Code     string  `json:"code"`
	Language string  `json:"language"`
	UserID   string  `json:"user_id"`
	LineItem string  `json:"lineitem"`
	IDToken  string  `json:"id_token"`
	MaxScore float64 `json:"max_score"`
}

// ExecuteResponse represents the response from code execution
type ExecuteResponse struct {
	Success bool                   `json:"success"`
	Result  *models.Judge0Response `json:"result,omitempty"`
	Score   float64                `json:"score,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// ExecuteHandler handles code execution and grade submission
func ExecuteHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("⚡ Code execution request received")

	// Parse request body
	var req ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Error parsing request: %v", err)
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Code == "" || req.Language == "" {
		log.Println("❌ Missing required fields: code or language")
		sendErrorResponse(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Load config
	cfg := config.LoadConfig()

	// Initialize services
	judge0Service := services.NewJudge0Service(cfg.Judge0AuthToken)
	agsService := services.NewAGSService(
		cfg.TokenEndpoint,
		cfg.ClientID,
		cfg.ClientSecret,
	)

	// Get language ID and submit code
	langID := config.GetLanguageID(req.Language)
	result, err := judge0Service.SubmitCode(req.Code, langID)
	if err != nil {
		log.Printf("❌ Judge0 error: %v", err)
		sendErrorResponse(w, fmt.Sprintf("Execution error: %v", err), http.StatusInternalServerError)
		return
	}

	// Calculate score
	score := judge0Service.CalculateScore(result, req.MaxScore)

	// Submit grade to Moodle if lineitem is available
	if req.LineItem != "" && req.UserID != "" {
		go func() {
			gradeReq := models.AGSGradeRequest{
				LineItemURL: req.LineItem,
				UserID:      req.UserID,
				Score:       score,
				MaxScore:    req.MaxScore,
				Comment:     fmt.Sprintf("Auto-graded at %s", time.Now().Format(time.RFC3339)),
			}

			if err := agsService.SubmitGrade(gradeReq); err != nil {
				log.Printf("⚠️ Failed to submit grade: %v", err)
			} else {
				log.Printf("✅ Grade submitted - User: %s, Score: %.2f/%.2f",
					req.UserID, score, req.MaxScore)
			}
		}()
	}

	// Send response
	response := ExecuteResponse{
		Success: true,
		Result:  result,
		Score:   score,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// // Helper function to calculate score based on execution result
// func calculateScore(result *Judge0Response, maxScore float64) float64 {
// 	if result == nil {
// 		return 0
// 	}

// 	// Simple scoring logic: if there's stdout and no errors, give max score
// 	if result.Stdout != nil && *result.Stdout != "" &&
// 		(result.Stderr == nil || *result.Stderr == "") &&
// 		(result.CompileOutput == nil || *result.CompileOutput == "") {
// 		return maxScore
// 	}

// 	return 0
// }

// Helper function to send error response
func sendErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ExecuteResponse{
		Success: false,
		Error:   message,
	})
}
