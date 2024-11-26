package main

import (
	"github.com/elmawardy/nutrix/backend/cmd"
	"github.com/elmawardy/nutrix/backend/common/config"
	"github.com/elmawardy/nutrix/backend/common/logger"
	"github.com/elmawardy/nutrix/backend/common/userio"
	"github.com/elmawardy/nutrix/backend/modules"
	"github.com/elmawardy/nutrix/backend/modules/core"
	"github.com/gorilla/mux"
)

func main() {

	logger := logger.NewZeroLog()
	conf := config.ConfigFactory("viper", "config.yaml", &logger)
	settings := config.Settings{}
	err := settings.LoadFromDB(conf)
	if err != nil {
		logger.Error(err.Error())
		panic("Can't load settings from DB")
	}

	logger.Info("Successfully connected to DB")

	prompter := &userio.BubbleTeaSeedablesPrompter{
		Logger: &logger,
	}

	router := mux.NewRouter()

	appmanager := modules.AppManager{
		Logger: &logger,
	}

	appmanager.LoadModule(&core.Core{
		Logger:   &logger,
		Config:   conf,
		Prompter: prompter,
		Settings: settings,
	}, "core").RegisterHttpHandlers(router).RegisterBackgroundWorkers().Save()

	appmanager.Ignite()

	modules, err := appmanager.GetModules()
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	rootCmd := cmd.RootProcess{
		Config:   conf,
		Logger:   &logger,
		Settings: settings,
		Router:   router,
		Modules:  modules,
		Prompter: prompter,
	}

	rootCmd.Execute()
}
