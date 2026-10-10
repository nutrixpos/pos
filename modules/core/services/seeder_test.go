package services

import (
	"context"
	"testing"

	"github.com/nutrixpos/pos/common/userio"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/require"
)

type fakePrompter struct{ confirm bool }

func (f fakePrompter) Confirmation(string) (bool, error) { return f.confirm, nil }
func (f fakePrompter) MultiChooseTree(string, []userio.PromptTreeElement) ([]userio.PromptTreeElement, error) {
	return nil, nil
}

func TestSeeder_SeedSettings_CreatesWhenMissing(t *testing.T) {
	env := newServiceEnv(t)
	svc := Seeder{Logger: env.Logger, Config: env.Config, Prompter: fakePrompter{confirm: true}}

	// Drop the settings created by the test harness so the insert branch runs.
	require.NoError(t, env.Client.Database(env.Config.Databases[0].Database).Collection("settings").Drop(context.Background()))
	require.NoError(t, svc.SeedSettings())
}

func TestSeeder_SeedProducts_ConfirmCreatesMaterials(t *testing.T) {
	env := newServiceEnv(t)
	svc := Seeder{Logger: env.Logger, Config: env.Config, Prompter: fakePrompter{confirm: true}}

	// No seeded materials: the prompter confirms and materials are created first.
	require.NoError(t, svc.SeedProducts())
}

func TestSeeder_SeedCategories_ConfirmCreatesProducts(t *testing.T) {
	env := newServiceEnv(t)
	svc := Seeder{Logger: env.Logger, Config: env.Config, Prompter: fakePrompter{confirm: true}}

	// Pre-create just the category so the "already exists" branch runs, then the
	// seeded products are missing and the prompter confirms creating them.
	seedServiceDoc(t, env, "categories", models.Category{Name: "CategorySeeded"})
	require.NoError(t, svc.SeedCategories())
}

func TestSeeder_SeedCategories_FreshConfirmTrue(t *testing.T) {
	env := newServiceEnv(t)
	svc := Seeder{Logger: env.Logger, Config: env.Config, Prompter: fakePrompter{confirm: true}}

	// Fresh env: the prompter confirms creating seeded products; the category
	// lookup then fails on the singular "ProductSeeded" name and returns an error.
	require.Error(t, svc.SeedCategories())
}

func TestSeeder_Seeds(t *testing.T) {
	env := newServiceEnv(t)
	svc := Seeder{Logger: env.Logger, Config: env.Config, Prompter: fakePrompter{confirm: true}}

	// Settings already exist (testutil seeds them): idempotent path.
	require.NoError(t, svc.SeedSettings())

	// Materials first so products can reference them.
	require.NoError(t, svc.SeedMaterials(true))

	// A second run with IsNewOnly skips existing materials.
	svc.IsNewOnly = true
	require.NoError(t, svc.SeedMaterials(true))

	require.NoError(t, svc.SeedProducts())
	require.NoError(t, svc.SeedCategories())

	// Running again with confirm=true exercises the "already exists" path.
	svc.IsNewOnly = false
	require.NoError(t, svc.SeedProducts())
	require.NoError(t, svc.SeedCategories())
	require.NoError(t, svc.SeedMaterials(false))
}
