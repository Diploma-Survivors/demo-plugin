package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-lti-provider/models"
	"net/http"
)

const (
	judge0URL = "http://localhost:2358/submissions"
)

// Judge0Service handles interaction with Judge0 API
type Judge0Service struct {
	BaseURL string
}

// NewJudge0Service creates a new Judge0Service instance
func NewJudge0Service(baseURL string) *Judge0Service {
	if baseURL == "" {
		baseURL = judge0URL
	}
	return &Judge0Service{BaseURL: baseURL}
}

// SubmitCode submits code to Judge0 for execution
func (s *Judge0Service) SubmitCode(code string, languageID int) (*models.Judge0Response, error) {
	submission := models.Submission{
		SourceCode: code,
		LanguageID: languageID,
	}

	jsonData, err := json.Marshal(submission)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal submission: %w", err)
	}

	// Submit with wait=true to get result immediately
	resp, err := http.Post(s.BaseURL+"?wait=true", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to POST to Judge0: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Judge0 returned status %d", resp.StatusCode)
	}

	var result models.Judge0Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Judge0 response: %w", err)
	}

	return &result, nil
}

// GetLanguageID returns the Judge0 language ID for a given language name
func (s *Judge0Service) GetLanguageID(language string) int {
	// Map of supported languages to Judge0 language IDs
	languageMap := map[string]int{
		"python": 71, // Python (3.8.1)
		"go":     60, // Go (1.13.5)
		"java":   62, // Java (OpenJDK 13.0.1)
		"cpp":    54, // C++ (GCC 9.2.0)
		"c":      50, // C (GCC 9.2.0)
		"nodejs": 63, // JavaScript (Node.js 12.14.0)
	}

	if id, ok := languageMap[language]; ok {
		return id
	}
	return 71 // Default to Python if language not found
}

// CalculateScore calculates the score based on execution result
func (s *Judge0Service) CalculateScore(result *models.Judge0Response, maxScore float64) float64 {
	if result == nil {
		return 0
	}

	// Simple scoring logic: if there's stdout and no errors, give max score
	if result.Stdout != nil && *result.Stdout != "" &&
		(result.Stderr == nil || *result.Stderr == "") &&
		(result.CompileOutput == nil || *result.CompileOutput == "") {
		return maxScore
	}

	return 0
}
