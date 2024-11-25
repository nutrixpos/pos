package config

import (
	"log"

	"github.com/elmawardy/nutrix/backend/common/logger"
	"github.com/spf13/viper"
)

func NewViperConfig(logger logger.ILogger) *ViperConfig {
	vc := new(ViperConfig)

	vc.v = viper.New()
	vc.logger = logger

	return vc
}

type ViperConfig struct {
	v      *viper.Viper
	logger logger.ILogger
}

func (vc *ViperConfig) ReadFile(path string) {
	vc.v.SetConfigFile(path)
}

func (vc *ViperConfig) GetConfig() (Config, error) {
	// Read the config file
	if err := vc.v.ReadInConfig(); err != nil {
		vc.logger.Error("Failed to read config file: %v", err)
		return Config{}, err
	}

	var config Config
	if err := vc.v.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return config, nil
}
