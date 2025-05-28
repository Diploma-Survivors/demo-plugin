package config

import (
	"os"
)

// LTI 1.3 Configuration
type Config struct {
	// Server settings
	Port string `env:"PORT" default:"8080"`

	// LTI Platform (Moodle) settings
	PlatformIssuer   string // Moodle's issuer URL
	PlatformJWKSURL  string // Moodle's JWKS endpoint
	PlatformTokenURL string // Moodle's OAuth2 token endpoint
	PlatformAuthURL  string // Moodle's OIDC auth endpoint

	// Tool settings
	ClientID     string // LTI Tool Client ID trong Moodle
	ClientSecret string // LTI Tool Client Secret (nếu cần)
	ToolIssuer   string // Tool's issuer URL (localhost cho dev)

	// Deployment settings
	DeploymentID string // LTI Deployment ID

	// AGS settings
	AGSScope string // Assignment and Grade Services scope

	// Judge0 settings
	Judge0URL       string
	Judge0AuthToken string // Nếu Judge0 có authentication

	// Security settings
	JWTSigningMethod string
	AllowedOrigins   []string

	// Moodle settings
	MoodleBaseURL  string `env:"MOODLE_BASE_URL" default:"http://localhost:8888"`
	AuthLoginURL   string `env:"MOODLE_AUTH_URL" default:"http://localhost:8888/mod/lti/auth.php"`
	TokenEndpoint  string `env:"MOODLE_TOKEN_URL" default:"http://localhost:8888/mod/lti/token.php"`
	KeysetEndpoint string `env:"MOODLE_KEYSET_URL" default:"http://localhost:8888/mod/lti/certs.php"`
	FrontendURL    string `env:"FRONTEND_URL" default:"http://localhost:3000"`
}

// LoadConfig loads configuration from environment variables với default values
func LoadConfig() *Config {
	return &Config{
		// Server settings
		Port: getEnv("PORT", "8080"),

		// Platform settings cho Moodle local
		PlatformIssuer:   getEnv("PLATFORM_ISSUER", "http://localhost:8888"),
		PlatformJWKSURL:  getEnv("PLATFORM_JWKS_URL", "http://localhost:8888/mod/lti/certs.php"),
		PlatformTokenURL: getEnv("PLATFORM_TOKEN_URL", "http://localhost:8888/mod/lti/token.php"),
		PlatformAuthURL:  getEnv("PLATFORM_AUTH_URL", "http://localhost:8888/mod/lti/auth.php"),

		// Tool settings
		ClientID:     getEnv("LTI_CLIENT_ID", "wAWXk7ifY0o9tCU"),
		ClientSecret: getEnv("LTI_CLIENT_SECRET", "your-client-secret"),
		ToolIssuer:   getEnv("TOOL_ISSUER", "http://localhost:8080"),

		// Deployment
		DeploymentID: getEnv("LTI_DEPLOYMENT_ID", "1"),

		// AGS
		AGSScope: getEnv("AGS_SCOPE", "https://purl.imsglobal.org/spec/lti-ags/scope/score"),

		// Judge0
		Judge0URL:       getEnv("JUDGE0_URL", "http://localhost:2358"),
		Judge0AuthToken: getEnv("JUDGE0_AUTH_TOKEN", ""),

		// Security
		JWTSigningMethod: getEnv("JWT_SIGNING_METHOD", "RS256"),
		AllowedOrigins:   []string{getEnv("ALLOWED_ORIGINS", "*")},

		// Moodle settings
		MoodleBaseURL:  getEnv("MOODLE_BASE_URL", "http://localhost:8888"),
		AuthLoginURL:   getEnv("MOODLE_AUTH_URL", "http://localhost:8888/mod/lti/auth.php"),
		TokenEndpoint:  getEnv("MOODLE_TOKEN_URL", "http://localhost:8888/mod/lti/token.php"),
		KeysetEndpoint: getEnv("MOODLE_KEYSET_URL", "http://localhost:8888/mod/lti/certs.php"),
		FrontendURL:    getEnv("FRONTEND_URL", "http://localhost:3000"),
	}
}

// getEnv gets environment variable với fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Validate kiểm tra cấu hình có hợp lệ không
func (c *Config) Validate() []string {
	var errors []string

	if c.PlatformIssuer == "http://localhost:8080" && c.ClientID == "wAWXk7ifY0o9tCU" {
		errors = append(errors, "⚠️ Using default values - configure LTI_CLIENT_ID for production")
	}

	if c.Port == "" {
		errors = append(errors, "PORT không được để trống")
	}

	return errors
}

// GetJudge0SubmissionURL returns full Judge0 submission URL
func (c *Config) GetJudge0SubmissionURL() string {
	return c.Judge0URL + "/submissions"
}

// GetToolLoginURL returns tool's login endpoint URL
func (c *Config) GetToolLoginURL() string {
	return c.ToolIssuer + "/lti/login"
}

// GetToolLaunchURL returns tool's launch endpoint URL
func (c *Config) GetToolLaunchURL() string {
	return c.ToolIssuer + "/lti/launch"
}

// GetToolJWKSURL returns tool's JWKS endpoint URL (nếu tool cũng cần serve JWKS)
func (c *Config) GetToolJWKSURL() string {
	return c.ToolIssuer + "/lti/jwks"
}

// Language mappings cho Judge0
var LanguageMap = map[string]int{
	"go":         75, // Go
	"python":     71, // Python 3
	"java":       62, // Java
	"javascript": 63, // JavaScript (Node.js)
	"cpp":        54, // C++
	"c":          50, // C
	"php":        68, // PHP
	"ruby":       72, // Ruby
	"rust":       73, // Rust
	"swift":      83, // Swift
}

// GetLanguageID returns Judge0 language ID from language name
func GetLanguageID(language string) int {
	if id, exists := LanguageMap[language]; exists {
		return id
	}
	return 75 // Default to Go
}
