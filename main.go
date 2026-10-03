// Package main is the entrypoint for the nutrix application.
//
// It sets up the config, logger, and userio, and then starts the main server.
//
// The main server is a gorilla/mux router that hosts the core, auth, and other
// modules.
package main

import (
	"os"

	"github.com/nutrixpos/pos/cmd"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/common/userio"
)

func main() {

	// Initialize the logger using ZeroLog
	logger := logger.NewZeroLog()

	// Create the configuration using the Viper config backend
	conf := config.ConfigFactory("viper", "config.yaml", &logger)

	// Initialize the prompter for user interaction
	prompter := &userio.BubbleTeaSeedablesPrompter{
		Logger: &logger,
	}
	// Initialize the root command process with configuration and modules
	rootCmd := cmd.RootProcess{
		Config:   conf,
		Logger:   &logger,
		Prompter: prompter,
	}

	// Execute the root command to start the application
	if err := rootCmd.Execute(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
