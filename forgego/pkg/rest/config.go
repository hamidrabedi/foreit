package rest

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Host string `mapstructure:"host" default:"0.0.0.0"`
	Port int    `mapstructure:"port" default:"8080" validate:"min=1,max=65535"`
}

// StaticConfig contains static file serving settings
type StaticConfig struct {
	Enable bool   `mapstructure:"enable" default:"true"`
	Path   string `mapstructure:"path" default:"/static"`
	Root   string `mapstructure:"root" default:"./public"`
}

