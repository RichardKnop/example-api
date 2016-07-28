package config

// DatabaseConfig stores database connection options
type DatabaseConfig struct {
	Type         string
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	MaxIdleConns int
	MaxOpenConns int
}

// OauthConfig stores oauth service configuration options
type OauthConfig struct {
	AccessTokenLifetime  int
	RefreshTokenLifetime int
	AuthCodeLifetime     int
}

// FacebookConfig stores Facebook app config
type FacebookConfig struct {
	AppID     string
	AppSecret string
}

// SendgridConfig stores sengrid configuration options
type SendgridConfig struct {
	APIKey string
}

// WebConfig stores web related config like scheme and host
type WebConfig struct {
	Scheme    string
	Host      string
	AppScheme string
	AppHost   string
}

// AppSpecificConfig stores app specific config
type AppSpecificConfig struct {
	PasswordResetLifetime int
	CompanyName           string
	CompanyEmail          string
}

// Config stores all configuration options
type Config struct {
	Database      DatabaseConfig
	Oauth         OauthConfig
	Facebook      FacebookConfig
	Sendgrid      SendgridConfig
	Web           WebConfig
	AppSpecific   AppSpecificConfig
	IsDevelopment bool
}
