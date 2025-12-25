package api

// APIConfig contains API module configuration
type APIConfig struct {
	Path string `mapstructure:"path" default:"/api/v1"`
}

