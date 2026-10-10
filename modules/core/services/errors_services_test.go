package services

import (
	"testing"

	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
)

func TestRecipeService_ErrorPaths(t *testing.T) {
	env := newServiceEnv(t)
	svc := newRecipeService(env)

	_, err := svc.GetProduct("missing")
	assert.Error(t, err)

	_, err = svc.GetReadyNumber("missing")
	assert.Error(t, err)

	_, err = svc.GetRecipeTree("missing")
	assert.Error(t, err)

	assert.Error(t, svc.ConsumeFromReady("missing", 1))

	_, _, err = svc.GetProducts(GetProductsParams{PageNumber: 1, PageSize: 10})
	assert.NoError(t, err)
}

func TestMaterialService_EditMissing(t *testing.T) {
	env := newServiceEnv(t)
	svc := newMaterialService(env, models.Settings{})
	assert.Error(t, svc.EditMaterial("missing", models.Material{Name: "x"}))
}

func TestSeeder_ExistingConfirmFalse(t *testing.T) {
	env := newServiceEnv(t)
	svc := Seeder{Logger: env.Logger, Config: env.Config, Prompter: fakePrompter{confirm: false}}

	// First call creates; second call finds existing data and the prompter
	// declines to reseed.
	assert.NoError(t, svc.SeedProducts())
	assert.NoError(t, svc.SeedProducts())

	assert.NoError(t, svc.SeedCategories())
	assert.NoError(t, svc.SeedCategories())

	// IsNewOnly short-circuits existing data.
	svc.IsNewOnly = true
	assert.NoError(t, svc.SeedProducts())
	assert.NoError(t, svc.SeedCategories())
}
