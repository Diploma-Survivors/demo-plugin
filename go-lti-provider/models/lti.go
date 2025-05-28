package models

// LTILaunchClaims represents the claims from LTI 1.3 launch JWT
type LTILaunchClaims struct {
	Issuer         string                 `json:"iss"`
	Subject        string                 `json:"sub"`
	Audience       []string               `json:"aud"`
	ExpirationTime int64                  `json:"exp"`
	IssuedAt       int64                  `json:"iat"`
	Nonce          string                 `json:"nonce"`
	MessageType    string                 `json:"https://purl.imsglobal.org/spec/lti/claim/message_type"`
	Version        string                 `json:"https://purl.imsglobal.org/spec/lti/claim/version"`
	DeploymentID   string                 `json:"https://purl.imsglobal.org/spec/lti/claim/deployment_id"`
	TargetLinkURI  string                 `json:"https://purl.imsglobal.org/spec/lti/claim/target_link_uri"`
	ResourceLink   ResourceLink           `json:"https://purl.imsglobal.org/spec/lti/claim/resource_link"`
	Context        Context                `json:"https://purl.imsglobal.org/spec/lti/claim/context"`
	Custom         map[string]interface{} `json:"https://purl.imsglobal.org/spec/lti/claim/custom,omitempty"`
	Roles          []string               `json:"https://purl.imsglobal.org/spec/lti/claim/roles"`
	EndpointClaim  *EndpointClaim         `json:"https://purl.imsglobal.org/spec/lti-ags/claim/endpoint,omitempty"`
}

type ResourceLink struct {
	ID          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type Context struct {
	ID    string   `json:"id"`
	Label string   `json:"label,omitempty"`
	Title string   `json:"title,omitempty"`
	Type  []string `json:"type,omitempty"`
}

type EndpointClaim struct {
	Scope    []string `json:"scope,omitempty"`
	LineItem string   `json:"lineitem,omitempty"`
}

// AGSGradeRequest represents a request to submit grade to Moodle
type AGSGradeRequest struct {
	LineItemURL string  `json:"lineitem_url"`
	UserID      string  `json:"user_id"`
	Score       float64 `json:"score"`
	MaxScore    float64 `json:"max_score"`
	Comment     string  `json:"comment"`
	AccessToken string  `json:"access_token,omitempty"`
}

// Grade represents the grade submission to Moodle AGS
type Grade struct {
	ScoreGiven       float64 `json:"scoreGiven"`
	ScoreMaximum     float64 `json:"scoreMaximum"`
	Comment          string  `json:"comment,omitempty"`
	ActivityProgress string  `json:"activityProgress"` // Initialized, InProgress, Submitted, Completed
	GradingProgress  string  `json:"gradingProgress"`  // FullyGraded, Pending, PendingManual, Failed, NotReady
	Timestamp        string  `json:"timestamp"`
	UserID           string  `json:"userId"`
}
