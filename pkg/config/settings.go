package config

// Settings represents framework settings structure
type Settings struct {
	Admin    AdminSettings
	App      AppSettings
	Security SecuritySettings
	Server   ServerSettings
	Database DatabaseSettings
}

// AppSettings contains application-level settings
type AppSettings struct {
	Name    string
	Env     string
	Version string
	Debug   bool
}

// ServerSettings contains HTTP server settings
type ServerSettings struct {
	Host         string
	Port         string
	ReadTimeout  int // seconds
	WriteTimeout int // seconds
}

// DatabaseSettings contains database connection settings
type DatabaseSettings struct {
	Driver          string
	Host            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	Port            int
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
	ConnMaxIdleTime int
}

// SecuritySettings contains security-related settings
type SecuritySettings struct {
	SecretKey     string
	CSRFSecretKey string
	SessionSecret string
}

// AdminSettings contains admin interface settings
type AdminSettings struct {
	Path        string
	Title       string
	HeaderTitle string
	SiteName    string
	Enabled     bool
}

// LoadSettings loads settings from config
func LoadSettings(cfg *Config) *Settings {
	return &Settings{
		App: AppSettings{
			Name:    cfg.GetString("app.name", "forge"),
			Env:     cfg.GetString("app.env", "development"),
			Debug:   cfg.GetBool("app.debug", true),
			Version: cfg.GetString("app.version", "0.1.0"),
		},
		Server: ServerSettings{
			Host:         cfg.GetString("server.host", "localhost"),
			Port:         cfg.GetString("server.port", "8000"),
			ReadTimeout:  cfg.GetInt("server.read_timeout", 30),
			WriteTimeout: cfg.GetInt("server.write_timeout", 30),
		},
		Database: DatabaseSettings{
			Driver:          cfg.GetString("database.driver", "postgres"),
			Host:            cfg.GetString("database.host", "localhost"),
			Port:            cfg.GetInt("database.port", 5432),
			User:            cfg.GetString("database.user", "postgres"),
			Password:        cfg.GetString("database.password", ""),
			Name:            cfg.GetString("database.name", "forge"),
			SSLMode:         cfg.GetString("database.sslmode", "disable"),
			MaxOpenConns:    cfg.GetInt("database.max_open_conns", 25),
			MaxIdleConns:    cfg.GetInt("database.max_idle_conns", 5),
			ConnMaxLifetime: cfg.GetInt("database.conn_max_lifetime", 300),
			ConnMaxIdleTime: cfg.GetInt("database.conn_max_idle_time", 600),
		},
		Security: SecuritySettings{
			SecretKey:     cfg.GetString("security.secret_key", "change-me-in-production"),
			CSRFSecretKey: cfg.GetString("security.csrf_secret_key", "change-me-in-production"),
			SessionSecret: cfg.GetString("security.session_secret", "change-me-in-production"),
		},
		Admin: AdminSettings{
			Enabled:     cfg.GetBool("admin.enabled", true),
			Path:        cfg.GetString("admin.path", "/admin"),
			Title:       cfg.GetString("admin.title", "forge Admin"),
			HeaderTitle: cfg.GetString("admin.header_title", "forge"),
			SiteName:    cfg.GetString("admin.site_name", "forge"),
		},
	}
}
