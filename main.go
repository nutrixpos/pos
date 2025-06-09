// Package main is the entrypoint for the nutrix application.
//
// It sets up the config, logger, and userio, and then starts the main server.
//
// The main server is a gorilla/mux router that hosts the core, auth, and other
// modules.
package main

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/cmd"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/common/userio"
	"github.com/nutrixpos/pos/modules"
	"github.com/nutrixpos/pos/modules/core"
	"github.com/nutrixpos/pos/modules/core/services"
	"github.com/nutrixpos/pos/modules/hubsync"
)

func main() {

	// Initialize the logger using ZeroLog
	logger := logger.NewZeroLog()

	// Create the configuration using the Viper config backend
	conf := config.ConfigFactory("viper", "config.yaml", &logger)

	seeder_svc := services.Seeder{
		Config: conf,
	}

	// make sure that settings bootstrapping data exists, it's idempotent
	err := seeder_svc.SeedSettings()
	if err != nil {
		panic(err)
	}

	settings_svc := services.SettingsService{
		Config: conf,
	}

	// Load settings from the database
	settings, err := settings_svc.GetSettings()

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

	appmanager.LoadModule(&hubsync.HubSyncModule{
		Logger: &logger,
		Config: conf,
	}, "hubsync").RegisterBackgroundWorkers().RegisterHttpHandlers(router).Save()

	// Ignite the app manager to start all modules
	appmanager.Run()

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
