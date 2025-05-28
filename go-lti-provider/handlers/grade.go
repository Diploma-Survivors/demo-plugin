package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// AGS Grade submission structures
type Grade struct {
	ScoreGiven       float64 `json:"scoreGiven"`
	ScoreMaximum     float64 `json:"scoreMaximum"`
	Comment          string  `json:"comment,omitempty"`
	ActivityProgress string  `json:"activityProgress"` // Initialized, InProgress, Submitted, Completed
	GradingProgress  string  `json:"gradingProgress"`  // FullyGraded, Pending, PendingManual, Failed, NotReady
	Timestamp        string  `json:"timestamp"`
	UserID           string  `json:"userId"`
}

// OAuth2 Token Response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// AGS Grade Request
type AGSGradeRequest struct {
	LineItemURL string  `json:"lineitem_url"`
	UserID      string  `json:"user_id"`
	Score       float64 `json:"score"`
	MaxScore    float64 `json:"max_score"`
	Comment     string  `json:"comment"`
	AccessToken string  `json:"access_token"`
}

// GradeHandler xử lý việc gửi điểm về Moodle qua AGS
func GradeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("📊 AGS Grade submission received")

	// Parse request body
	var gradeReq AGSGradeRequest
	if err := json.NewDecoder(r.Body).Decode(&gradeReq); err != nil {
		log.Printf("❌ Error parsing grade request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if gradeReq.LineItemURL == "" || gradeReq.UserID == "" {
		log.Println("❌ Missing required fields: lineitem_url or user_id")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Get access token if not provided
	accessToken := gradeReq.AccessToken
	if accessToken == "" {
		token, err := getAccessToken()
		if err != nil {
			log.Printf("❌ Failed to get access token: %v", err)
			http.Error(w, "Failed to authenticate with platform", http.StatusInternalServerError)
			return
		}
		accessToken = token
	}

	// Submit grade to Moodle
	err := submitGradeToMoodle(gradeReq, accessToken)
	if err != nil {
		log.Printf("❌ Failed to submit grade: %v", err)
		http.Error(w, fmt.Sprintf("Failed to submit grade: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Grade submitted successfully - User: %s, Score: %.2f/%.2f",
		gradeReq.UserID, gradeReq.Score, gradeReq.MaxScore)

	// Return success response
	response := map[string]interface{}{
		"success": true,
		"message": "Grade submitted successfully",
		"data": map[string]interface{}{
			"user_id":   gradeReq.UserID,
			"score":     gradeReq.Score,
			"max_score": gradeReq.MaxScore,
			"comment":   gradeReq.Comment,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getAccessToken() (string, error) {
	// OAuth2 Client Credentials flow cho Moodle local
	tokenURL := "http://localhost:8888/mod/lti/token.php"
	clientID := "wAWXk7ifY0o9tCU"
	clientSecret := "your-client-secret"
	scope := "https://purl.imsglobal.org/spec/lti-ags/scope/score"

	// Prepare form data
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

func submitGradeToMoodle(gradeReq AGSGradeRequest, accessToken string) error {
	// Create AGS Grade payload
	grade := Grade{
		ScoreGiven:       gradeReq.Score,
		ScoreMaximum:     gradeReq.MaxScore,
		Comment:          gradeReq.Comment,
		ActivityProgress: "Completed",
		GradingProgress:  "FullyGraded",
		Timestamp:        time.Now().Format(time.RFC3339),
		UserID:           gradeReq.UserID,
	}

	jsonData, err := json.Marshal(grade)
	if err != nil {
		return fmt.Errorf("failed to marshal grade data: %w", err)
	}

	// AGS Score submission endpoint
	// URL format: {line_item_url}/scores
	scoreURL := gradeReq.LineItemURL + "/scores"

	req, err := http.NewRequest("POST", scoreURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create grade request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/vnd.ims.lis.v1.score+json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to submit grade: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("grade submission failed with status %d", resp.StatusCode)
	}

	log.Printf("✅ Grade submitted to AGS endpoint: %s", scoreURL)
	return nil
}

// Helper function để test grade submission từ launch
func SubmitTestGrade(lineItemURL, userID string, score, maxScore float64) error {
	gradeReq := AGSGradeRequest{
		LineItemURL: lineItemURL,
		UserID:      userID,
		Score:       score,
		MaxScore:    maxScore,
		Comment:     "Auto-graded by LTI Tool",
	}

	accessToken, err := getAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	return submitGradeToMoodle(gradeReq, accessToken)
}

// ParseScore parses score từ string và validate
func ParseScore(scoreStr string, maxScore float64) (float64, error) {
	if scoreStr == "" {
		return 0, fmt.Errorf("empty score")
	}

	score, err := strconv.ParseFloat(scoreStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid score format: %w", err)
	}

	if score < 0 {
		return 0, fmt.Errorf("score cannot be negative")
	}

	if score > maxScore {
		return maxScore, nil // Cap at maximum score
	}

	return score, nil
}
