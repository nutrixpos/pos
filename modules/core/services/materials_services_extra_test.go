package services

import (
	"testing"

	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaterialService_AddEditDelete(t *testing.T) {
	env := newServiceEnv(t)
	svc := newMaterialService(env, models.Settings{})

	require.NoError(t, svc.AddComponent(models.Material{Name: "Flour", Unit: "kg"}, "u"))

	materials, err := svc.GetMaterials(1, 10)
	require.NoError(t, err)
	require.Len(t, materials, 1)
	id := materials[0].Id

	require.NoError(t, svc.EditMaterial(id, models.Material{Name: "Flour2", Unit: "g", Settings: models.MaterialSettings{StockAlertTreshold: 5}}))
	got := getMaterial(t, env, id)
	assert.Equal(t, "Flour2", got.Name)
	assert.Equal(t, "g", got.Unit)
	assert.Equal(t, 5.0, got.Settings.StockAlertTreshold)

	require.NoError(t, svc.DeleteMaterial(id))
	_, err = svc.GetMaterial(id)
	assert.Error(t, err)
}

func TestMaterialService_CalculateCosts(t *testing.T) {
	env := newServiceEnv(t)
	insertMaterial(t, env, models.Material{
		Id: "mat-calc",
		Entries: []models.MaterialEntry{
			{Id: "e1", Quantity: 10, PurchaseQuantity: 10, PurchasePrice: 100},
			{Id: "e2", Quantity: 10, PurchaseQuantity: 10, PurchasePrice: 200},
		},
	})
	svc := newMaterialService(env, models.Settings{})

	avg, err := svc.CalculateMaterialAverageCost("mat-calc", 2)
	require.NoError(t, err)
	assert.Equal(t, 30.0, avg)

	exact, err := svc.CalculateMaterialExactCost("e1", "mat-calc", 2)
	require.NoError(t, err)
	assert.Equal(t, 20.0, exact)

	_, err = svc.CalculateMaterialExactCost("nope", "mat-calc", 2)
	assert.Error(t, err)
}

func TestMaterialService_GetMaterialEntries(t *testing.T) {
	env := newServiceEnv(t)
	insertMaterial(t, env, models.Material{
		Id: "mat-entries",
		Entries: []models.MaterialEntry{
			{Id: "e1", Quantity: 1},
			{Id: "e2", Quantity: 2},
			{Id: "e3", Quantity: 3},
		},
	})
	svc := newMaterialService(env, models.Settings{})

	entries, total, err := svc.GetMaterialEntries("mat-entries", GetMaterialEntriesParams{PageNumber: 0, PageSize: 2})
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, entries, 2)

	_, _, err = svc.GetMaterialEntries("mat-entries", GetMaterialEntriesParams{PageNumber: -1})
	assert.Error(t, err)

	// Entries beyond the last page return an empty slice.
	entries, _, err = svc.GetMaterialEntries("mat-entries", GetMaterialEntriesParams{PageNumber: 100, PageSize: 2})
	require.NoError(t, err)
	assert.Empty(t, entries)
}
