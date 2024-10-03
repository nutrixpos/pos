package main

import (
	"log"
	"net/http"
	"time"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
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
		panic("Can't load settings from DB")
	}

	router := mux.NewRouter()

	modules_manager := modules.ModulesManager{}
	modules_manager.RegisterModule("core", &logger, core.NewBuilder(conf, settings), router)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
