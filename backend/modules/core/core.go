package core

import (
	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/common/userio"
	"github.com/elmawardy/nutrix/modules"
	auth_mw "github.com/elmawardy/nutrix/modules/auth/middlewares"
	"github.com/elmawardy/nutrix/modules/core/handlers"
	"github.com/elmawardy/nutrix/modules/core/middlewares"
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

func (cmb *CoreModuleBuilder) RegisterHttpHandlers(router *mux.Router, prefix string) modules.IModuleBuilder {

	auth_svc := auth_mw.NewZitadelAuth(cmb.Config)

	router.Handle(prefix+"/api/sales_logs", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSalesLog(cmb.Config, cmb.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/salesperday", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSalesPerDay(cmb.Config, cmb.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/entry", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteEntry(cmb.Config, cmb.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materialentry", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.PushMaterialEntry(cmb.Config, cmb.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/material", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.AddMaterial(cmb.Config, cmb.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/order", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetOrder(cmb.Config, cmb.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materiallogs", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetMaterialLogs(cmb.Config, cmb.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materials", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetMaterials(cmb.Config, cmb.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materialcost", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CalculateMaterialCost(cmb.Config, cmb.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/categories", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetCategories(cmb.Config, cmb.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/startorder", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.StartOrder(cmb.Config, cmb.Logger, cmb.Settings), "admin", "chef"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetOrders(cmb.Config, cmb.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/orderstash", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.OrderStash(cmb.Config, cmb.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orderremovestash", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.OrderRemoveFromStash(cmb.Config, cmb.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/ordergetstashed", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetStashedOrders(cmb.Config, cmb.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/submitorder", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.SubmitOrder(cmb.Config, cmb.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/finishorder", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.FinishOrder(cmb.Config, cmb.Logger), "admin", "chef"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/recipeavailability", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetRecipeAvailability(cmb.Config, cmb.Logger), "admin", "chef", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/recipetree", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetRecipeTree(cmb.Config, cmb.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/products", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetProducts(cmb.Config, cmb.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/addproduct", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.InesrtNewProduct(cmb.Config, cmb.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/deleteproduct", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteProduct(cmb.Config, cmb.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/updateproduct", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateProduct(cmb.Config, cmb.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/productgetready", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetProductReadyNumber(cmb.Config, cmb.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/editmaterial", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.EditMaterial(cmb.Config, cmb.Logger), "admin"))).Methods("POST", "OPTIONS")

	notification_service, err := services.SpawnNotificationService("melody", cmb.Logger, cmb.Config)
	if err != nil {
		cmb.Logger.Error(err.Error())
		panic(err)
	}

	router.Handle(prefix+"/ws", handlers.HandleNotificationsWsRequest(cmb.Config, cmb.Logger, notification_service))

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
