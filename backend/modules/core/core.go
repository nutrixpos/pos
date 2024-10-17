package core

import (
	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/common/userio"
	"github.com/elmawardy/nutrix/modules"
	"github.com/elmawardy/nutrix/modules/core/handlers"
	"github.com/elmawardy/nutrix/modules/core/services"
	"github.com/gorilla/mux"
)

func NewBuilder(config config.Config, settings config.Settings) *CoreModuleBuilder {
	cmb := new(CoreModuleBuilder)
	cmb.Config = config
	cmb.Settings = settings

	return cmb
}

type Core struct {
	Logger   logger.ILogger
	Config   config.Config
	Settings config.Settings
	Prompter userio.Prompter
}

type CoreModuleBuilder struct {
	Logger   logger.ILogger
	Config   config.Config
	Settings config.Settings
	Prompter userio.Prompter
}

func (cmb *CoreModuleBuilder) SetLogger(logger logger.ILogger) modules.IModuleBuilder {
	cmb.Logger = logger
	return cmb
}

func (cmb *CoreModuleBuilder) SetPrompter(prompter userio.Prompter) modules.IModuleBuilder {
	cmb.Prompter = prompter
	return cmb
}

func (c *Core) Seed(entities []string, is_new_only bool) error {

	seedService := services.Seeder{
		Logger:    c.Logger,
		Config:    c.Config,
		Prompter:  c.Prompter,
		IsNewOnly: is_new_only,
	}

	seedablesMap := make(map[string]bool, len(entities))

	for index := range entities {
		seedablesMap[entities[index]] = true
	}

	if _, ok := seedablesMap["materials"]; ok {
		c.Logger.Info("seeding materials ...")
		err := seedService.SeedMaterials(true)
		if err != nil {
			c.Logger.Error(err.Error())
			return err
		}
	}

	if _, ok := seedablesMap["products"]; ok {
		c.Logger.Info("seeding products ...")
		err := seedService.SeedProducts()
		if err != nil {
			c.Logger.Error(err.Error())
			return err
		}
	}

	if _, ok := seedablesMap["categories"]; ok {
		c.Logger.Info("seeding categories ...")
		err := seedService.SeedCategories()
		if err != nil {
			c.Logger.Error(err.Error())
			return err
		}
	}

	return nil
}

func (c *Core) GetSeedables() (entities []string, err error) {
	c.Logger.Info("Getting seedables...")

	return []string{
		"products",
		"materials",
		"materialentries",
		"categories",
		"settings",
	}, nil
}

func (cmb *CoreModuleBuilder) RegisterHttpHandlers(router *mux.Router) modules.IModuleBuilder {

	router.Handle("/api/sales_logs", handlers.GetSalesLog(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/entry", handlers.DeleteEntry(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/materialentry", handlers.PushMaterialEntry(cmb.Config, cmb.Logger)).Methods("POST", "OPTIONS")
	router.Handle("/api/material", handlers.AddMaterial(cmb.Config, cmb.Logger)).Methods("POST", "OPTIONS")
	router.Handle("/api/order", handlers.GetOrder(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/materiallogs", handlers.GetMaterialLogs(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/materials", handlers.GetMaterials(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/materialcost", handlers.CalculateMaterialCost(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/categories", handlers.GetCategories(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/startorder", handlers.StartOrder(cmb.Config, cmb.Logger, cmb.Settings)).Methods("POST", "OPTIONS")
	router.Handle("/api/orders", handlers.GetOrders(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/orderstash", handlers.OrderStash(cmb.Config, cmb.Logger)).Methods("POST", "OPTIONS")
	router.Handle("/api/orderremovestash", handlers.OrderRemoveFromStash(cmb.Config, cmb.Logger)).Methods("POST", "OPTIONS")
	router.Handle("/api/ordergetstashed", handlers.GetStashedOrders(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/submitorder", handlers.SubmitOrder(cmb.Config, cmb.Logger)).Methods("POST", "OPTIONS")
	router.Handle("/api/finishorder", handlers.FinishOrder(cmb.Config, cmb.Logger)).Methods("POST", "OPTIONS")
	router.Handle("/api/recipeavailability", handlers.GetRecipeAvailability(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/recipetree", handlers.GetRecipeTree(cmb.Config, cmb.Logger)).Methods("GET")
	router.Handle("/api/editmaterial", handlers.EditMaterial(cmb.Config, cmb.Logger)).Methods("POST", "OPTIONS")

	notification_service, err := services.SpawnNotificationService("melody", cmb.Logger, cmb.Config)
	if err != nil {
		cmb.Logger.Error(err.Error())
		panic(err)
	}

	router.Handle("/ws", handlers.HandleNotificationsWsRequest(cmb.Config, cmb.Logger, notification_service))

	return cmb
}

func (cmb *CoreModuleBuilder) Build() modules.BaseModule {
	return &Core{
		Logger:   cmb.Logger,
		Config:   cmb.Config,
		Prompter: cmb.Prompter,
		Settings: cmb.Settings,
	}

}
