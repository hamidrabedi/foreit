package i18n

// I18nConfig contains internationalization settings
type I18nConfig struct {
	Enable       bool   `mapstructure:"enable" default:"false"`
	LocalesPath  string `mapstructure:"locales_path" default:"./locales"`
	DefaultLocale string `mapstructure:"default_locale" default:"en"`
}

