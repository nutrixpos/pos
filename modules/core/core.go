// Package core contains the core module of the nutrix application.
//
// The core module contains the main business logic of the application, including
// the data models, services, and handlers for the core features of the
// application.
package core

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/helpers"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/common/userio"
	"github.com/nutrixpos/pos/modules"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	"github.com/nutrixpos/pos/modules/core/handlers"
	"github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/nutrixpos/pos/modules/core/services"
)

// Core is the main struct for the core module.
//
// It contains the necessary fields for the core module to function, including
// the logger, config, settings, prompter, and notification service.
type Core struct {
	// Logger is the logger object for the core module.
	Logger logger.ILogger

	// Config is the config object for the core module.
	Config config.Config

	// Settings is the settings object for the core module.
	Settings models.Settings

	// Prompter is the prompter object for the core module.
	Prompter userio.Prompter

	// NotificationSvc is the notification service object for the core module.
	NotificationSvc services.INotificationService
}

// OnStart is called when the core module is started.
func (c *Core) OnStart() func() error {
	return func() error {
		return nil
	}
}

// OnEnd is called when the core module is ended.
func (c *Core) OnEnd() func() {
	return func() {

	}
}

// Seed seeds the database with sample data.
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

// GetSeedables returns a list of seedables.
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

// RegisterBackgroundWorkers registers background workers.
func (c *Core) RegisterBackgroundWorkers() []modules.Worker {

	if c.NotificationSvc == nil {
		notification_service, err := services.SpawnNotificationSingletonSvc("melody", c.Logger, c.Config)
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

// RegisterHttpHandlers registers HTTP handlers.
func (c *Core) RegisterHttpHandlers(router *mux.Router, prefix string) {

	auth_svc := auth_mw.NewZitadelAuth(c.Config)

	c.Logger.Info("Successfully conntected to Zitadel")

	router.Handle(prefix+"/api/customers/{id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateCustomer(c.Config, c.Logger), "admin", "cashier"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/customers/{id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteCustomer(c.Config, c.Logger, c.Settings), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/customers/{id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetCustomer(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/customers", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetCustomers(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/customers", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.AddCustomer(c.Config, c.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/logs/salesperday", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSalesPerDay(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/logs/salesperday/exportcsv", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.ExportSalesCSV(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materials", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetMaterials(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materials", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.AddMaterial(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/materials/{id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.EditMaterial(c.Config, c.Logger), "admin"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/materials/{id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteMaterial(c.Config, c.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/materials/{id}/logs", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetMaterialLogs(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materials/{material_id}/entries/{entry_id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteEntry(c.Config, c.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/materials/{material_id}/entries/{entry_id}/cost", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CalculateMaterialExactCost(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materials/{id}/entries", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.PushMaterialEntry(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/materials/{material_id}/entries", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetMaterialEntries(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materials/{material_id}/avgcost", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CalculateMaterialAverageCost(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/categories", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetCategories(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/categories", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.InsertCategory(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/categories/{id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteCategory(c.Config, c.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/categories/{id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateCategory(c.Config, c.Logger), "admin"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/orders", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetOrders(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/orders", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.SubmitOrder(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetOrder(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteOrder(c.Config, c.Logger), "admin", "cashier"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}/start", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.StartOrder(c.Config, c.Logger, c.Settings), "admin", "chef"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{order_id}/logs", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetOrderLogs(c.Config, c.Logger, c.Settings), "admin", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}/cancel", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CancelOrder(c.Config, c.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}/finish", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.FinishOrder(c.Config, c.Logger, c.Settings), "admin", "chef"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}/pay", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.Payorder(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}/printkitchenreceipt", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.PrintKitchenReceipt(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}/printclientreceipt", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.PrintClientReceipt(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{order_id}/addtips", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.OrderAddTip(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/orders/{order_id}/removetips", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.OrderRemoveTip(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/orders/{order_id}/items/{item_id}/refund", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.RefundOrderItem(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{order_id}/items/{item_id}/waste", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.WasteOrderItem(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/products/availability", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetRecipeAvailability(c.Config, c.Logger), "admin", "chef", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/products/{id}/recipetree", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetRecipeTree(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/products/{id}/image", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateProductImage(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/products/{id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetProduct(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/products/{id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteProduct(c.Config, c.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/products/{id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateProduct(c.Config, c.Logger), "admin"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/products", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetProducts(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/products", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.InesrtNewProduct(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/settings", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSettings(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/settings", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateSettings(c.Config, c.Logger), "admin"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/languages", middlewares.AllowCors(handlers.GetAvailableLanguages(c.Config, c.Logger))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/languages/{code}", middlewares.AllowCors(handlers.GetLanguage(c.Config, c.Logger))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/disposals/{id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetDisposal(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/disposals/{id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteProduct(c.Config, c.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/disposals/{id}", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateProduct(c.Config, c.Logger), "admin"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/disposals", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetDisposals(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/disposals", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.InesrtNewProduct(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")

	if c.NotificationSvc == nil {
		notification_service, err := services.SpawnNotificationSingletonSvc("melody", c.Logger, c.Config)
		if err != nil {
			c.Logger.Error(err.Error())
			panic(err)
		}
		c.NotificationSvc = notification_service
	}

	// create public folder if doesn't exist in the uploads dir directory
	publicPath := helpers.ResolveOsEnvPath(c.Config.UploadsPath)
	if _, err := os.Stat(publicPath); os.IsNotExist(err) {
		err = os.MkdirAll(publicPath, 0755)
		if err != nil {
			c.Logger.Error(err.Error())
			panic(err)
		}
	}

	// Serve static files from the "static" directory
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir(helpers.ResolveOsEnvPath(c.Config.UploadsPath)))))

	router.Handle(prefix+"/ws", handlers.HandleNotificationsWsRequest(c.Config, c.Logger, c.NotificationSvc))
}
