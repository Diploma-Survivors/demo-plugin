package models

// Submission represents a code submission to Judge0
type Submission struct {
	SourceCode string `json:"source_code"`
	LanguageID int    `json:"language_id"`
	Stdin      string `json:"stdin,omitempty"`
}

// Judge0Response represents the response from Judge0 API
type Judge0Response struct {
	Token         string  `json:"token,omitempty"`
	Status        Status  `json:"status,omitempty"`
	Stdout        *string `json:"stdout"`
	Stderr        *string `json:"stderr"`
	CompileOutput *string `json:"compile_output"`
	Message       *string `json:"message"`
	ExitCode      *int    `json:"exit_code"`
	ExitSignal    *int    `json:"exit_signal"`
	Time          *string `json:"time"`
	Memory        *int    `json:"memory"`
}

// Status represents the execution status from Judge0
type Status struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

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
	Success bool            `json:"success"`
	Result  *Judge0Response `json:"result,omitempty"`
	Score   float64         `json:"score,omitempty"`
	Error   string          `json:"error,omitempty"`
}
