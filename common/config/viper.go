// Package config implements the config interface for the viper config backend.
package config

import (
	"log"
	"strings"

	"github.com/nutrixpos/pos/common/logger"
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

	tables := make(map[string]string)

	for k, v := range vc.v.GetStringMap("databases.0.tables") {
		tables[k] = v.(string)
	}

	databases := make([]Database, 1)
	databases[0] = Database{
		Host:     vc.v.GetString("databases.0.host"),
		Port:     vc.v.GetInt("databases.0.port"),
		Database: vc.v.GetString("databases.0.database"),
		Username: vc.v.GetString("databases.0.username"),
		Password: vc.v.GetString("databases.0.password"),
		Type:     vc.v.GetString("databases.0.type"),
		Name:     vc.v.GetString("databases.0.name"),
		Tables:   tables,
	}

	zitadel_domain := vc.v.GetString("zitadel.domain")
	zitadel_port := vc.v.GetInt("zitadel.port")

	// Unmarshal the config into the config struct.
	var config Config
	if err := vc.v.Unmarshal(&config); err != nil {
		// Log the error.
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	config.Databases = databases
	config.Zitadel.Domain = zitadel_domain
	config.Zitadel.Port = zitadel_port

	return config, nil
}
