package middleware

// SecurityConfig contains security-related settings
type SecurityConfig struct {
	SecretKey string `mapstructure:"secret_key" validate:"required"`
	Debug     bool   `mapstructure:"debug" default:"false"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level  string `mapstructure:"level" default:"info"` // "debug", "info", "warn", "error"
	Format string `mapstructure:"format" default:"json"` // "json", "text"
}

