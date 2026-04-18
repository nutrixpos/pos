// Package cmd handles the command-line interface and operations for the application.
package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/helpers"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/common/userio"
	"github.com/nutrixpos/pos/modules/core"
	"github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/nutrixpos/pos/modules/core/services"
	"github.com/nutrixpos/pos/modules/hubsync"
	"gopkg.in/yaml.v2"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules"
	"github.com/spf13/cobra"
)

// RootProcess holds the root process configuration.
type RootProcess struct {
	Config   config.Config
	Logger   logger.ILogger
	Settings models.Settings
	cmd      *cobra.Command
	Router   *mux.Router
	Modules  map[string]modules.IBaseModule
	Prompter userio.Prompter
}

// Execute starts the root process.
func (root *RootProcess) Execute() error {

	cobra.MousetrapHelpText = ""
	setupRun := false

	root.cmd = &cobra.Command{
		Use:   "nutrix",
		Short: "The next level restaurant management system",
		Long: `Nutrix is an open source project aiming to make restaurant management a seamless and efficient experience.
		Free forever and distributed under the GPT-2 license. https://github.com/nutrixpos/pos`,

		Run: func(cmd *cobra.Command, args []string) {

			frontendServerStartChan := make(chan bool, 1)
			frontendOpenBrowserChan := make(chan bool, 1)

			if root.Config.ServeFrontEnd {
				go startFrontendServer(frontendServerStartChan, root.Logger)
				go func(config config.Config, logger logger.ILogger, startChan chan bool) {

					<-startChan

					logger.Info("Opening browser window")
					time.Sleep(2 * time.Second)
					err := helpers.OpenURL("http://localhost:8080")
					if err != nil {
						logger.Error(err.Error())
					}
				}(root.Config, root.Logger, frontendOpenBrowserChan)
			}

			// Create a new HTTP router
			root.Router = mux.NewRouter()

			root.Router.Handle("/api/setup/status", middlewares.AllowCors(
				func() http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/json")

						if root.Config.Databases[0].Host == "" {
							w.Write([]byte(`{"setup":false}`))
							return
						}

						w.Write([]byte(`{"setup":true}`))
						return
					}
				}(),
			)).Methods("GET", "OPTIONS")

			// If no database host is configured, serve a setup endpoint so the user
			// can provide connection details via the browser (Setup.vue). The process
			// exits with code 0 after writing the config so the process manager can
			// restart the app with the new configuration.
			if len(root.Config.Databases) == 0 || root.Config.Databases[0].Host == "" {

				setupRun = true
				stopChan := make(chan struct{})

				root.Logger.Info("No database host configured — starting setup server on :8000")

				root.Router.Handle("/api/setup/test-connection", middlewares.AllowCors(
					func() http.HandlerFunc {
						return func(w http.ResponseWriter, r *http.Request) {
							w.Header().Set("Content-Type", "application/json")
							w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

							if r.Method == http.MethodOptions {
								w.WriteHeader(http.StatusOK)
								return
							}

							if r.Method != http.MethodPost {
								http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
								return
							}

							var body struct {
								Host     string `json:"host"`
								Port     int    `json:"port"`
								Database string `json:"database"`
								Username string `json:"username"`
								Password string `json:"password"`
							}

							if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
								http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
								return
							}

							if body.Host == "" || body.Port == 0 || body.Database == "" {
								http.Error(w, "host, port and database are required", http.StatusBadRequest)
								return
							}

							testCfg := &config.Config{
								Databases: []config.Database{
									{
										Host:     body.Host,
										Port:     body.Port,
										Database: body.Database,
										Username: body.Username,
										Password: body.Password,
									},
								},
							}

							client, err := common.GetDatabaseClient(root.Logger, testCfg)
							if err != nil {
								http.Error(w, "Connection failed: "+err.Error(), http.StatusBadRequest)
								return
							}

							ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
							defer cancel()

							err = client.Ping(ctx, nil)
							if err != nil {
								http.Error(w, "Ping failed: "+err.Error(), http.StatusBadRequest)
								return
							}

							w.WriteHeader(http.StatusOK)
							w.Write([]byte(`{"success":true}`))
						}
					}(),
				)).Methods("POST", "OPTIONS")

				root.Router.Handle("/api/setup/config", middlewares.AllowCors(
					func() http.HandlerFunc {
						return func(w http.ResponseWriter, r *http.Request) {
							// CORS headers for the Vue dev-server proxy
							w.Header().Set("Content-Type", "application/json")
							w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

							if r.Method == http.MethodOptions {
								w.WriteHeader(http.StatusOK)
								return
							}

							if r.Method != http.MethodPost {
								http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
								return
							}

							var body struct {
								Host     string `json:"host"`
								Port     int    `json:"port"`
								Database string `json:"database"`
								Username string `json:"username"`
								Password string `json:"password"`
							}

							if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
								http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
								return
							}

							if body.Host == "" || body.Port == 0 || body.Database == "" {
								http.Error(w, "host, port and database are required", http.StatusBadRequest)
								return
							}

							// Patch (or create) the first database entry.
							databases := root.Config.Databases
							if len(databases) == 0 {
								databases = []config.Database{{}}
							}
							core_db := &databases[0]
							if core_db == nil {
								core_db = &config.Database{}
							}
							core_db.Host = body.Host
							core_db.Port = body.Port
							core_db.Database = body.Database
							core_db.Username = body.Username
							core_db.Password = body.Password
							core_db.Type = "mongo"
							databases[0] = *core_db

							// Marshal back to YAML and write.
							updated, err := yaml.Marshal(root.Config)
							if err != nil {
								http.Error(w, "cannot marshal config.yaml: "+err.Error(), http.StatusInternalServerError)
								return
							}

							if err := os.WriteFile("config.yaml", updated, 0644); err != nil {
								http.Error(w, "cannot write config.yaml: "+err.Error(), http.StatusInternalServerError)
								return
							}

							root.Logger.Info(fmt.Sprintf("Setup complete — database host set to %s:%d/%s. Awaiting context exit...", body.Host, body.Port, body.Database))

							w.WriteHeader(http.StatusOK)
							close(stopChan)
						}
					}(),
				)).Methods("POST", "OPTIONS")

				srv := &http.Server{
					Handler: root.Router,
					Addr:    "0.0.0.0:8000",
					// Good practice: enforce timeouts for servers you create!
					WriteTimeout: 15 * time.Second,
					ReadTimeout:  15 * time.Second,
				}

				frontendServerStartChan <- true
				frontendOpenBrowserChan <- true

				listener, err := net.Listen("tcp4", srv.Addr)
				if err != nil {
					log.Fatal(err)
				}

				go func() {
					if err := srv.Serve(listener); err != nil {
						root.Logger.Error("setup server error: " + err.Error())

						if !errors.Is(err, http.ErrServerClosed) {
							os.Exit(1)
						}
					}
				}()

				<-stopChan // Block until /quit is called

				// Gracefully shut down with a timeout
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := srv.Shutdown(ctx); err != nil {
					log.Fatalf("Server Shutdown Failed:%+v", err)
				}
				fmt.Println("Setup server exited gracefully")
			}

			root.Router.Handle("/api/setup/config", http.NotFoundHandler())
			root.Router.Handle("/api/setup/test-connection", http.NotFoundHandler())

			seeder_svc := services.Seeder{
				Config: root.Config,
				Logger: root.Logger,
			}

			// make sure that settings bootstrapping data exists, it's idempotent
			err := seeder_svc.SeedSettings()
			if err != nil {
				panic(err)
			}

			settings_svc := services.SettingsService{
				Config: root.Config,
			}

			// Load settings from the database
			settings, err := settings_svc.GetSettings()

			if err != nil {
				// Log and panic if settings can't be loaded
				root.Logger.Error(err.Error())
				panic("Can't load settings from DB")
			}

			// Log successful database connection
			root.Logger.Info("Successfully connected to DB")

			// Initialize the app manager with logger
			appmanager := modules.AppManager{
				Logger: root.Logger,
			}

			// Load the core module, register HTTP handlers and background workers, and save the module
			appmanager.LoadModule(&core.Core{
				Logger:   root.Logger,
				Config:   root.Config,
				Prompter: root.Prompter,
				Settings: settings,
			}, "core").RegisterHttpHandlers(root.Router).RegisterBackgroundWorkers().Save()

			appmanager.LoadModule(&hubsync.HubSyncModule{
				Logger: root.Logger,
				Config: root.Config,
			}, "hubsync").RegisterBackgroundWorkers().RegisterHttpHandlers(root.Router).Save()

			// Ignite the app manager to start all modules
			appmanager.Run()

			// Retrieve all registered modules
			modules, err := appmanager.GetModules()
			if err != nil {
				// Log error and panic if modules can't be retrieved
				root.Logger.Error(err.Error())
				panic(err)
			}

			root.Modules = modules

			srv := &http.Server{
				Handler: root.Router,
				Addr:    "0.0.0.0:8000",
				// Good practice: enforce timeouts for servers you create!
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			if !setupRun {
				frontendServerStartChan <- true
				frontendOpenBrowserChan <- true
			}

			listener, err := net.Listen("tcp4", srv.Addr)
			if err != nil {
				log.Fatal(err)
			}

			log.Fatal(srv.Serve(listener))
		},
	}

	seedService := SeedProcess{
		Config:   root.Config,
		Logger:   root.Logger,
		Settings: root.Settings,
		Router:   root.Router,
		Modules:  root.Modules,
	}

	seedCmd, err := seedService.GetCmd(root.Prompter)
	if err != nil {
		return err
	}

	root.cmd.AddCommand(seedCmd)

	if err := root.cmd.Execute(); err != nil {
		return err
	}

	return nil
}

func startFrontendServer(startChan chan bool, logger logger.ILogger) {

	<-startChan

	// Create a new ServeMux for the static file server to keep its handlers separate
	// from the main API router.
	staticMux := http.NewServeMux()

	// Create a FileServer to serve files from a "web" directory.
	// You should create a folder named `web` and place your static files there.
	fs := http.FileServer(http.Dir("./mnt/frontend"))

	// Handle requests for static files. The http.StripPrefix ensures that
	// a request for /static/index.html looks for /index.html in the `web` directory.
	staticMux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Optionally, serve the index.html from the root path of the static server.
	staticMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./mnt/frontend/index.html")
	})

	srv := &http.Server{
		Handler:      staticMux,
		Addr:         "0.0.0.0:8080", // Listen on a new, separate port
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("Serving static files on http://localhost:8080/")
	log.Fatal(srv.ListenAndServe())
}
