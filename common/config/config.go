package config

import (
	"github.com/elmawardy/nutrix/common/logger"
)

// ConfigFactory creates a Config object based on the provided type and path
func ConfigFactory(t string, path string, logger logger.ILogger) Config {
	switch t {
	case "viper":
		viper_config := NewViperConfig(logger)
		viper_config.ReadFile(path)

		config, err := viper_config.GetConfig()
		if err != nil {
			logger.Error("can't reat config")
		}

		return config
	}

	return Config{}
}

// IConfig interface defines methods for config management
type IConfig interface {
	ReadFile(path string)
	BindEnv() error
	GetConfig() (config Config)
}

// ZitadelConfig holds the configuration for Zitadel
type ZitadelConfig struct {
	Domain  string `mapstructure:"domain"`
	Port    int    `mapstructure:"port"`
	KeyPath string `mapstructure:"key_path"`
}

// Config represents the overall configuration structure
type Config struct {
	Databases   []Database    `mapstructure:"databases"`
	Zitadel     ZitadelConfig `mapstructure:"zitadel"`
	Env         string        `mapstructure:"env"`
	TimeZone    string        `mapstructure:"timezone"`
	UploadsPath string        `mapstructure:"uploads_path"`
}

// Database holds the configuration for database connections
type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Type     string `mapstructure:"type"`
	Name     string `mapstructure:"name"`
	Database string `mapstructure:"database"`
}
