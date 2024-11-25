package cmd

import (
	"fmt"

	"github.com/elmawardy/nutrix/backend/common/config"
	"github.com/elmawardy/nutrix/backend/common/logger"
	"github.com/elmawardy/nutrix/backend/common/userio"
	"github.com/elmawardy/nutrix/backend/modules"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

type SeedProcess struct {
	Config    config.Config
	Logger    logger.ILogger
	Settings  config.Settings
	Router    *mux.Router
	Modules   map[string]modules.IBaseModule
	IsNewOnly bool
}

func (sp *SeedProcess) GetCmd(prompter userio.Prompter) (*cobra.Command, error) {

	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed db with data for dev/test purposes.",
		Run: func(cmd *cobra.Command, args []string) {
			sp.Logger.Info("Seeding..")
			err := sp.Seed(sp.Modules, prompter)
			if err != nil {
				sp.Logger.Error(err.Error())
				panic(err)
			}
		},
	}

	cmd.PersistentFlags().BoolVar(&sp.IsNewOnly, "new-only", false, "seed only non existing data, and leave old data untouched")

	return cmd, nil

}

func (sp *SeedProcess) Seed(mods map[string]modules.IBaseModule, prompter userio.Prompter) error {
	sp.Logger.Info("Seeding database...")
	seedable_modules_prompt_elements := []userio.PromptTreeElement{}
	seedableModules := map[string]modules.ISeederModule{}

	for moduleName, module := range mods {
		if seederModule, ok := module.(modules.ISeederModule); ok {

			seedable_module_prompt_element := userio.PromptTreeElement{
				Title:       moduleName,
				Selected:    false,
				SubElements: []userio.PromptTreeElement{},
			}

			sp.Logger.Info(fmt.Sprintf("Seeding module: %s", moduleName))

			s, err := seederModule.GetSeedables()
			if err != nil {
				return err
			}

			for _, seedable := range s {

				sub_element := userio.PromptTreeElement{
					Title:    seedable,
					Selected: false,
				}

				seedable_module_prompt_element.SubElements = append(seedable_module_prompt_element.SubElements, sub_element)
			}

			seedable_modules_prompt_elements = append(seedable_modules_prompt_elements, seedable_module_prompt_element)
			seedableModules[moduleName] = seederModule
		}
	}

	selected, err := prompter.MultiChooseTree("\nChoose which collections to seed from which modules:\n\n", seedable_modules_prompt_elements)
	if err != nil {
		return err
	}

	for _, selectedSeedableModule := range selected {
		if selectedSeedableModule.Selected {
			sp.Logger.Info(fmt.Sprintf("Seeding: %s", selectedSeedableModule.Title))

			selected_module_seedables := []string{}

			for _, selectedselectedSeedable := range selectedSeedableModule.SubElements {
				if selectedselectedSeedable.Selected {
					selected_module_seedables = append(selected_module_seedables, selectedselectedSeedable.Title)
				}
			}

			seedableModules[selectedSeedableModule.Title].Seed(selected_module_seedables, sp.IsNewOnly)
		}
	}

	return nil
}
