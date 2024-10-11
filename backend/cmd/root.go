package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/common/userio"
	"github.com/elmawardy/nutrix/modules"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

type RootProcess struct {
	Config   config.Config
	Logger   logger.ILogger
	Settings config.Settings
	cmd      *cobra.Command
	Router   *mux.Router
	Modules  map[string]modules.BaseModule
	Prompter userio.Prompter
}

func (root *RootProcess) Execute() error {

	root.cmd = &cobra.Command{
		Use:   "nutrix",
		Short: "The next level restaurant management system",
		Long: `Nutrix is an open source project aiming to make restaurant management a seamless and efficient experience.
		Free forever and distributed under the MIT license. https://github.com/elmawardy/nutrix`,

		Run: func(cmd *cobra.Command, args []string) {
			srv := &http.Server{
				Handler: root.Router,
				Addr:    "127.0.0.1:8000",
				// Good practice: enforce timeouts for servers you create!
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			log.Fatal(srv.ListenAndServe())
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
