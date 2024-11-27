// Package main is the entrypoint for the nutrix application.
//
// It sets up the config, logger, and userio, and then starts the main server.
//
// The main server is a gorilla/mux router that hosts the core, auth, and other
// modules.
package main

import (
	"github.com/elmawardy/nutrix/cmd"
	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/common/userio"
	"github.com/elmawardy/nutrix/modules"
	"github.com/elmawardy/nutrix/modules/core"
	"github.com/gorilla/mux"
)

func main() {

	// Initialize the logger using ZeroLog
	logger := logger.NewZeroLog()

	// Create the configuration using the Viper config backend
	conf := config.ConfigFactory("viper", "config.yaml", &logger)

	// Initialize settings structure
	settings := config.Settings{}

	// Load settings from the database
	err := settings.LoadFromDB(conf)
	if err != nil {
		// Log and panic if settings can't be loaded
		logger.Error(err.Error())
		panic("Can't load settings from DB")
	}

	// Log successful database connection
	logger.Info("Successfully connected to DB")

	// Initialize the prompter for user interaction
	prompter := &userio.BubbleTeaSeedablesPrompter{
		Logger: &logger,
	}

	// Create a new HTTP router
	router := mux.NewRouter()

	// Initialize the app manager with logger
	appmanager := modules.AppManager{
		Logger: &logger,
	}

	// Load the core module, register HTTP handlers and background workers, and save the module
	appmanager.LoadModule(&core.Core{
		Logger:   &logger,
		Config:   conf,
		Prompter: prompter,
		Settings: settings,
	}, "core").RegisterHttpHandlers(router).RegisterBackgroundWorkers().Save()

	// Ignite the app manager to start all modules
	appmanager.Ignite()

	// Retrieve all registered modules
	modules, err := appmanager.GetModules()
	if err != nil {
		// Log error and panic if modules can't be retrieved
		logger.Error(err.Error())
		panic(err)
	}

	// Initialize the root command process with configuration and modules
	rootCmd := cmd.RootProcess{
		Config:   conf,
		Logger:   &logger,
		Settings: settings,
		Router:   router,
		Modules:  modules,
		Prompter: prompter,
	}

	// Execute the root command to start the application
	rootCmd.Execute()
}
