package mypackage

// MyPackageConfig contains configuration for a third-party package
// This demonstrates how external packages can integrate with Forge
type MyPackageConfig struct {
	Enable bool   `mapstructure:"enable" default:"true"`
	Path   string `mapstructure:"path" default:"/mypackage"`
	APIKey string `mapstructure:"api_key" validate:"required"`
}

