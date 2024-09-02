package core

import (
	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules"
	"github.com/elmawardy/nutrix/modules/core/handlers"
	"github.com/gorilla/mux"
)

func NewBuilder(config config.Config) *CoreModuleBuilder {
	cmb := new(CoreModuleBuilder)
	cmb.Config = config

	return cmb
}

type Core struct {
	Logger logger.ILogger
	Config config.Config
}

type CoreModuleBuilder struct {
	Logger logger.ILogger
	Config config.Config
}

func (cmb *CoreModuleBuilder) SetLogger(logger logger.ILogger) modules.IModuleBuilder {
	cmb.Logger = logger
	return cmb
}

func (cmb *CoreModuleBuilder) RegisterHttpHandlers(router *mux.Router) modules.IModuleBuilder {

	router.Handle("/api/sales_logs", handlers.GetSalesLog(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/entry", handlers.DeleteEntry(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/componententry", handlers.PushComponentEntry(cmb.Config, cmb.Logger)).Methods("POST", "OPTIONS")
	router.Handle("/api/component", handlers.AddComponent(cmb.Config, cmb.Logger)).Methods("POST", "OPTIONS")
	router.Handle("/api/order", handlers.GetOrder(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/componentlogs", handlers.GetComponentLogs(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/components", handlers.GetComponents(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/categories", handlers.GetCategories(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/startorder", handlers.StartOrder(cmb.Config, cmb.Logger)).Methods("POST", "OPTIONS")
	router.Handle("/api/startorder2", handlers.StartOrder2(cmb.Config, cmb.Logger)).Methods("POST", "OPTIONS")
	router.Handle("/api/orders", handlers.GetOrders(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/submitorder", handlers.SubmitOrder(cmb.Config, cmb.Logger)).Methods("POST", "OPTIONS")
	router.Handle("/api/finishorder", handlers.FinishOrder(cmb.Config, cmb.Logger)).Methods("POST", "OPTIONS")
	router.Handle("/api/finishorder2", handlers.FinishOrder2(cmb.Config, cmb.Logger)).Methods("POST", "OPTIONS")
	router.Handle("/api/recipeavailability", handlers.GetRecipeAvailability(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/recipetree", handlers.GetRecipeTree(cmb.Config, cmb.Logger)).Methods("GET")

	return cmb
}

func (cmb *CoreModuleBuilder) Build() modules.IModule {
	return &Core{Logger: cmb.Logger, Config: cmb.Config}
}
