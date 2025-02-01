// Package config implements the config interface for the viper config backend.
package config

import (
	"log"
	"strings"

	"github.com/elmawardy/nutrix/common/logger"
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

func (vc *ViperConfig) BindAllEnv() (err error) {

	// Bind all environment variables with a prefix of the config file name
	// and substitute . with _.
	vc.v.AutomaticEnv()
	vc.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// if err = vc.v.BindEnv("databases.0.name", "DATABASES_0_NAME"); err != nil {
	// 	return err
	// }

	// if err = vc.v.BindEnv("databases.0.type", "DATABASES_0_TYPE"); err != nil {
	// 	return err
	// }

	// if err = vc.v.BindEnv("databases.0.host", "DATABASES_0_HOST"); err != nil {
	// 	return err
	// }

	// if err = vc.v.BindEnv("databases.0.port", "DATABASES_0_PORT"); err != nil {
	// 	return err
	// }

	// if err = vc.v.BindEnv("databases.0.database", "DATABASES_0_DATABASE"); err != nil {
	// 	return err
	// }

	// if err = vc.v.BindEnv("databases.0.username", "DATABASES_0_USERNAME"); err != nil {
	// 	return err
	// }

	// if err = vc.v.BindEnv("databases.0.password", "DATABASES_0_PASSWORD"); err != nil {
	// 	return err
	// }

	return nil
}

// GetConfig reads the config file and unmarshals it into the config struct.
func (vc *ViperConfig) GetConfig() (Config, error) {

	// Read the config file.
	if err := vc.v.ReadInConfig(); err != nil {
		// Log the error.
		vc.logger.Error("Failed to read config file: %v", err)
		return Config{}, err
	}

	vc.BindAllEnv()

	databases := make([]Database, 1)
	databases[0] = Database{
		Host:     vc.v.GetString("databases.0.host"),
		Port:     vc.v.GetInt("databases.0.port"),
		Database: vc.v.GetString("databases.0.database"),
		Username: vc.v.GetString("databases.0.username"),
		Password: vc.v.GetString("databases.0.password"),
		Type:     vc.v.GetString("databases.0.type"),
		Name:     vc.v.GetString("databases.0.name"),
	}
	zitadel_domain := vc.v.GetString("zitadel.domain")

	// Unmarshal the config into the config struct.
	var config Config
	if err := vc.v.Unmarshal(&config); err != nil {
		// Log the error.
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	config.Databases = databases
	config.Zitadel.Domain = zitadel_domain

	return config, nil
}
