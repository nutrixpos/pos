// Package cmd handles the command-line interface and operations for the application.
package cmd

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/modules/core/models"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/common/userio"
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

	root.cmd = &cobra.Command{
		Use:   "nutrix",
		Short: "The next level restaurant management system",
		Long: `Nutrix is an open source project aiming to make restaurant management a seamless and efficient experience.
		Free forever and distributed under the GPT-2 license. https://github.com/nutrixpos/pos`,

		Run: func(cmd *cobra.Command, args []string) {
			srv := &http.Server{
				Handler: root.Router,
				Addr:    "0.0.0.0:8000",
				// Good practice: enforce timeouts for servers you create!
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
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
