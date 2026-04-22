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
	auth_handlers "github.com/nutrixpos/pos/modules/auth/handlers"
	"github.com/nutrixpos/pos/modules/core/handlers"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/nutrixpos/pos/modules/core/services"
	"github.com/nutrixpos/pos/common"
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

	var auth_svc auth_mw.IAuthService

	if c.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(c.Config.Auth.JWTSecret, c.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(c.Config, jwtUtil)
		c.Logger.Info("Using internal JWT authentication")
	} else if c.Config.Zitadel.Enabled {
		auth_svc = auth_mw.NewZitadelAuth(c.Config)
		c.Logger.Info("Using Zitadel authentication")
	} else {
		auth_svc = auth_mw.NewNoAuth(c.Config)
		c.Logger.Info("Authentication disabled")
	}

	if c.Config.Auth.Enabled {
		client, err := common.GetDatabaseClient(c.Logger, &c.Config)
		if err != nil {
			c.Logger.Error("failed to get database client", "error", err)
			return
		}
		usersCollection := client.Database(c.Config.Databases[0].Database).Collection("users")
		authHandler := auth_handlers.NewAuthHandler(c.Config, c.Logger, usersCollection)

		router.Handle(prefix+"/api/auth/login", core_middlewares.AllowCors(http.HandlerFunc(authHandler.Login))).Methods("POST", "OPTIONS")
		router.Handle(prefix+"/api/auth/register", core_middlewares.AllowCors(http.HandlerFunc(authHandler.Register))).Methods("POST", "OPTIONS")
		router.Handle(prefix+"/api/auth/me", core_middlewares.AllowCors(auth_svc.AllowAuthenticated(http.HandlerFunc(authHandler.Me)))).Methods("GET", "OPTIONS")
		router.Handle(prefix+"/api/auth/users", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(http.HandlerFunc(authHandler.GetUsers), "superuser"))).Methods("GET", "OPTIONS")
		router.Handle(prefix+"/api/auth/users", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(http.HandlerFunc(authHandler.DeleteUser), "superuser"))).Methods("DELETE", "OPTIONS")
	}

	router.Handle(prefix+"/api/customers/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateCustomer(c.Config, c.Logger), "admin", "cashier"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/customers/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteCustomer(c.Config, c.Logger, c.Settings), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/customers/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetCustomer(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/customers", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetCustomers(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/customers", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.AddCustomer(c.Config, c.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/logs/salesperday", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSalesPerDay(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/logs/salesperday/exportcsv", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.ExportSalesCSV(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materials", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetMaterials(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materials", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.AddMaterial(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/materials/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.EditMaterial(c.Config, c.Logger), "admin"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/materials/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteMaterial(c.Config, c.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/materials/{id}/logs", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetMaterialLogs(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materials/{material_id}/entries/{entry_id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteEntry(c.Config, c.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/materials/{material_id}/entries/{entry_id}/cost", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CalculateMaterialExactCost(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materials/{id}/entries", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.PushMaterialEntry(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/materials/{material_id}/entries", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetMaterialEntries(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materials/{material_id}/avgcost", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CalculateMaterialAverageCost(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/categories", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetCategories(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/categories", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.InsertCategory(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/categories/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteCategory(c.Config, c.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/categories/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateCategory(c.Config, c.Logger), "admin"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/orders", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetOrders(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/orders", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.SubmitOrder(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetOrder(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteOrder(c.Config, c.Logger), "admin", "cashier"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}/start", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.StartOrder(c.Config, c.Logger, c.Settings), "admin", "chef"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{order_id}/logs", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetOrderLogs(c.Config, c.Logger, c.Settings), "admin", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}/cancel", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CancelOrder(c.Config, c.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}/finish", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.FinishOrder(c.Config, c.Logger, c.Settings), "admin", "chef"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}/pay", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.Payorder(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}/printkitchenreceipt", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.PrintKitchenReceipt(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}/printclientreceipt", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.PrintClientReceipt(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{order_id}/addtips", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.OrderAddTip(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/orders/{order_id}/removetips", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.OrderRemoveTip(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/orders/{order_id}/items/{item_id}/refund", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.RefundOrderItem(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{order_id}/items/{item_id}/waste", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.WasteOrderItem(c.Config, c.Logger, c.Settings), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/products/availability", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetRecipeAvailability(c.Config, c.Logger), "admin", "chef", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/products/{id}/recipetree", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetRecipeTree(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/products/{id}/image", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateProductImage(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/products/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetProduct(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/products/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteProduct(c.Config, c.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/products/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateProduct(c.Config, c.Logger), "admin"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/products", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetProducts(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/products", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.InesrtNewProduct(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/settings", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSettings(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/settings", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateSettings(c.Config, c.Logger), "admin"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/languages", core_middlewares.AllowCors(handlers.GetAvailableLanguages(c.Config, c.Logger))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/languages/{code}", core_middlewares.AllowCors(handlers.GetLanguage(c.Config, c.Logger))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/disposals/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetDisposal(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/disposals/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteProduct(c.Config, c.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/disposals/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateProduct(c.Config, c.Logger), "admin"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/disposals", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetDisposals(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/disposals", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.InesrtNewProduct(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")

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
