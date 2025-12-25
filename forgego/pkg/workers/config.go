package workers

// WorkersConfig contains workers module configuration
type WorkersConfig struct {
	Enable      bool   `mapstructure:"enable" default:"true"`
	RedisAddr   string `mapstructure:"redis_addr" default:"localhost:6379"`
	RedisPassword string `mapstructure:"redis_password" default:""`
	RedisDB     int    `mapstructure:"redis_db" default:"0"`
	Concurrency int    `mapstructure:"concurrency" default:"10" validate:"min=1"`
}

