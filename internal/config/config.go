package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type JWTConfig struct {
	Secret        string `mapstructure:"secret"`
	ExpirySeconds int    `mapstructure:"expiry_seconds"`
}

type AppConfig struct {
	Env      string `mapstructure:"env"`
	Port     string `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
	Timezone string `mapstructure:"timezone"`
}

func Load() (Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("configs")
	v.AddConfigPath(".")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("app.env", "development")
	v.SetDefault("app.port", "8080")
	v.SetDefault("app.log_level", "info")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.name", "ela")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.timezone", "UTC")
	v.SetDefault("jwt.secret", "change-me-in-production")
	v.SetDefault("jwt.expiry_seconds", 3600)

	if err := v.ReadInConfig(); err != nil {
		_, notFound := err.(viper.ConfigFileNotFoundError)
		if !notFound {
			return Config{}, fmt.Errorf("read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}
