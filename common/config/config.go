package config

import (
	"github.com/nutrixpos/pos/common/logger"
)

// ConfigFactory creates a Config object based on the provided type and path
func ConfigFactory(t string, path string, logger logger.ILogger) Config {
	switch t {
	case "viper":
		viper_config := NewViperConfig(logger)
		viper_config.ReadFile(path)

		config, err := viper_config.GetConfig()
		if err != nil {
			logger.Error("can't read config")
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
	Domain  string `mapstructure:"domain" yaml:"domain"`
	Port    int    `mapstructure:"port" yaml:"port"`
	KeyPath string `mapstructure:"key_path" yaml:"key_path"`
	Enabled bool   `mapstructure:"enabled" yaml:"enabled"`
}

// AuthConfig holds the configuration for internal authentication
type AuthConfig struct {
	JWTSecret    string `mapstructure:"jwt_secret" yaml:"jwt_secret"`
	JWTExpireHrs int    `mapstructure:"jwt_expire_hrs" yaml:"jwt_expire_hrs"`
	Enabled      bool   `mapstructure:"enabled" yaml:"enabled"`
}

// Config represents the overall configuration structure
type Config struct {
	Databases     []Database    `mapstructure:"databases" yaml:"databases"`
	Auth          AuthConfig    `mapstructure:"auth" yaml:"auth"`
	Zitadel       ZitadelConfig `mapstructure:"zitadel" yaml:"zitadel"`
	Env           string        `mapstructure:"env" yaml:"env"`
	TimeZone      string        `mapstructure:"timezone" yaml:"timezone"`
	UploadsPath   string        `mapstructure:"uploads_path" yaml:"uploads_path"`
	ServeFrontEnd bool          `mapstructure:"serve_frontend" yaml:"serve_frontend"`
}

// Database holds the configuration for database connections
type Database struct {
	Host     string            `mapstructure:"host" yaml:"host"`
	Port     int               `mapstructure:"port" yaml:"port"`
	Username string            `mapstructure:"username" yaml:"username"`
	Password string            `mapstructure:"password" yaml:"password"`
	Type     string            `mapstructure:"type" yaml:"type"`
	Name     string            `mapstructure:"name" yaml:"name"`
	Database string            `mapstructure:"database" yaml:"database"`
	Tables   map[string]string `mapstructure:"tables" yaml:"tables"`
}
