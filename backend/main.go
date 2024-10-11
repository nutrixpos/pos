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

	logger := logger.NewZeroLog()
	conf := config.ConfigFactory("viper", "config.yaml", &logger)
	settings := config.Settings{}
	err := settings.LoadFromDB(conf)
	if err != nil {
		logger.Error(err.Error())
		panic("Can't load settings from DB")
	}

	prompter := &userio.BubbleTeaSeedablesPrompter{
		Logger: &logger,
	}

	router := mux.NewRouter()

	modules_manager := modules.ModulesManager{}
	modules_manager.RegisterModule("core", &logger, core.NewBuilder(conf, settings), prompter, router)

	modules, err := modules_manager.GetModules()
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	rootCmd := cmd.RootProcess{
		Config:   conf,
		Logger:   &logger,
		Settings: settings,
		Modules:  modules,
		Prompter: prompter,
	}

	rootCmd.Execute()
}
