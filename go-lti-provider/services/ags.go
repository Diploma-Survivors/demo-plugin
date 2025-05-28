package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-lti-provider/models"
	"net/http"
	"time"
)

// AGSService handles interaction with Moodle's Assignment and Grade Services
type AGSService struct {
	TokenURL     string
	ClientID     string
	ClientSecret string
}

// NewAGSService creates a new AGSService instance
func NewAGSService(tokenURL, clientID, clientSecret string) *AGSService {
	return &AGSService{
		TokenURL:     tokenURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

// SubmitGrade submits a grade to Moodle AGS
func (s *AGSService) SubmitGrade(req models.AGSGradeRequest) error {
	// Get access token if not provided
	accessToken := req.AccessToken
	if accessToken == "" {
		token, err := s.getAccessToken()
		if err != nil {
			return fmt.Errorf("failed to get access token: %w", err)
		}
		accessToken = token
	}

	// Create grade payload
	grade := models.Grade{
		ScoreGiven:       req.Score,
		ScoreMaximum:     req.MaxScore,
		Comment:          req.Comment,
		ActivityProgress: "Completed",
		GradingProgress:  "FullyGraded",
		Timestamp:        time.Now().Format(time.RFC3339),
		UserID:           req.UserID,
	}

	jsonData, err := json.Marshal(grade)
	if err != nil {
		return fmt.Errorf("failed to marshal grade data: %w", err)
	}

	// Submit grade to AGS endpoint
	scoreURL := req.LineItemURL + "/scores"
	httpReq, err := http.NewRequest("POST", scoreURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create grade request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/vnd.ims.lis.v1.score+json")
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to submit grade: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("grade submission failed with status %d", resp.StatusCode)
	}

	return nil
}

// getAccessToken gets an OAuth2 access token from Moodle
func (s *AGSService) getAccessToken() (string, error) {
	data := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s&scope=%s",
		s.ClientID,
		s.ClientSecret,
		"https://purl.imsglobal.org/spec/lti-ags/scope/score")

	req, err := http.NewRequest("POST", s.TokenURL, bytes.NewBufferString(data))
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

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}
