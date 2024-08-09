package config

import "github.com/elmawardy/nutrix/common/logger"

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

type IConfig interface {
	ReadFile(path string)
	GetConfig() (config Config)
}

type Config struct {
	Databases    []Database
	Env          string `mapstructure:"env"`
	JwtSecretKey string `mapstructure:"jwt_secret_key"`
	TimeZone     string `mapstructure:"timezone"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Type     string `mapstructure:"type"`
	Name     string `mapstructure:"name"`
	Database string `mapstructure:"database"`
}
