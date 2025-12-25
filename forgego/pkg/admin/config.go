package admin

// AdminConfig contains admin module configuration
type AdminConfig struct {
	Enable       bool   `mapstructure:"enable" default:"true"`
	Path         string `mapstructure:"path" default:"/admin"`
	TemplatePath string `mapstructure:"template_path" default:"./templates/admin"`
}

