package core

import (
	"time"

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

type Core struct {
	Logger          logger.ILogger
	Config          config.Config
	Settings        config.Settings
	Prompter        userio.Prompter
	NotificationSvc services.INotificationService
}

func (c *Core) OnStart() func() {
	return func() {

	}
}

func (c *Core) OnEnd() func() {
	return func() {

	}
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

func (c *Core) RegisterBackgroundWorkers() []modules.Worker {

	if c.NotificationSvc == nil {
		notification_service, err := services.SpawnNotificationService("melody", c.Logger, c.Config)
		if err != nil {
			c.Logger.Error(err.Error())
			panic(err)
		}
		c.NotificationSvc = notification_service
	}

	workers := []modules.Worker{
		{
			Interval: 1 * time.Hour,
			Task: func() {
				services.CheckExpirationDates(c.Logger, c.Config, c.NotificationSvc)
			},
		},
	}

	return workers
}

func (c *Core) RegisterHttpHandlers(router *mux.Router, prefix string) {

	auth_svc := auth_mw.NewZitadelAuth(c.Config)

	router.Handle(prefix+"/api/sales_logs", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSalesLog(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/salesperday", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSalesPerDay(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/entry", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteEntry(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materialentry", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.PushMaterialEntry(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/material", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.AddMaterial(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/order", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetOrder(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materiallogs", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetMaterialLogs(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materials", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetMaterials(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materialcost", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CalculateMaterialCost(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/categories", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetCategories(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/updatecategory", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateCategory(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/deletecategory", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteCategory(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/insertcategory", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.InsertCategory(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/startorder", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.StartOrder(c.Config, c.Logger, c.Settings), "admin", "chef"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetOrders(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/orderstash", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.OrderStash(c.Config, c.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orderremovestash", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.OrderRemoveFromStash(c.Config, c.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/ordergetstashed", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetStashedOrders(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/ordercancel", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CancelOrder(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/submitorder", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.SubmitOrder(c.Config, c.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/finishorder", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.FinishOrder(c.Config, c.Logger), "admin", "chef"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/recipeavailability", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetRecipeAvailability(c.Config, c.Logger), "admin", "chef", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/recipetree", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetRecipeTree(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/products", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetProducts(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/addproduct", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.InesrtNewProduct(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/deleteproduct", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteProduct(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/updateproduct", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateProduct(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/productgetready", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetProductReadyNumber(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/editmaterial", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.EditMaterial(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")

	if c.NotificationSvc == nil {
		notification_service, err := services.SpawnNotificationService("melody", c.Logger, c.Config)
		if err != nil {
			c.Logger.Error(err.Error())
			panic(err)
		}
		c.NotificationSvc = notification_service
	}

	router.Handle(prefix+"/ws", handlers.HandleNotificationsWsRequest(c.Config, c.Logger, c.NotificationSvc))
}
