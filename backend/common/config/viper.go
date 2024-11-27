// Package config implements the config interface for the viper config backend.
package config

import (
	"log"

	"github.com/elmawardy/nutrix/backend/common/logger"
	"github.com/spf13/viper"
)

// NewViperConfig creates a new viper config backend.
func NewViperConfig(logger logger.ILogger) *ViperConfig {
	vc := new(ViperConfig)

	// The viper instance.
	vc.v = viper.New()
	// The logger.
	vc.logger = logger

	return vc
}

// ViperConfig is the viper config backend.
type ViperConfig struct {
	// The viper instance.
	v *viper.Viper
	// The logger.
	logger logger.ILogger
}

// ReadFile reads the config file from the path.
func (vc *ViperConfig) ReadFile(path string) {
	// Set the config file for viper.
	vc.v.SetConfigFile(path)
}

// GetConfig reads the config file and unmarshals it into the config struct.
func (vc *ViperConfig) GetConfig() (Config, error) {
	// Read the config file.
	if err := vc.v.ReadInConfig(); err != nil {
		// Log the error.
		vc.logger.Error("Failed to read config file: %v", err)
		return Config{}, err
	}

	// Unmarshal the config into the config struct.
	var config Config
	if err := vc.v.Unmarshal(&config); err != nil {
		// Log the error.
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return config, nil
}
