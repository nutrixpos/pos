package services

import (
	"testing"

	"github.com/nutrixpos/pos/common/customerrors"
	"github.com/nutrixpos/pos/internal/testutil"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRecipeService(env *testutil.TestEnv) *RecipeService {
	return &RecipeService{Logger: env.Logger, Config: env.Config}
}

func TestRecipeService_CRUD(t *testing.T) {
	env := newServiceEnv(t)
	svc := newRecipeService(env)

	created, err := svc.InsertNew(models.Product{Name: "Pizza", Price: 12, Quantity: 1})
	require.NoError(t, err)
	require.NotEmpty(t, created.Id)

	got, err := svc.GetProduct(created.Id)
	require.NoError(t, err)
	assert.Equal(t, "Pizza", got.Name)

	require.NoError(t, svc.UpdateProduct(created.Id, models.Product{Name: "Pizza XL", Price: 15}))
	got, err = svc.GetProduct(created.Id)
	require.NoError(t, err)
	assert.Equal(t, "Pizza XL", got.Name)
	assert.Equal(t, 15.0, got.Price)

	products, total, err := svc.GetProducts(GetProductsParams{PageNumber: 1, PageSize: 10, Search: "Pizza"})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, products, 1)

	require.NoError(t, svc.DeleteProduct(created.Id))
	_, err = svc.GetProduct(created.Id)
	assert.Error(t, err)
}

func TestRecipeService_ReadyAndConsume(t *testing.T) {
	env := newServiceEnv(t)
	svc := newRecipeService(env)

	created, err := svc.InsertNew(models.Product{Name: "Pizza", Price: 12, Ready: 5})
	require.NoError(t, err)

	ready, err := svc.GetReadyNumber(created.Id)
	require.NoError(t, err)
	assert.Equal(t, 5.0, ready)

	require.NoError(t, svc.ConsumeFromReady(created.Id, 2))
	ready, err = svc.GetReadyNumber(created.Id)
	require.NoError(t, err)
	assert.Equal(t, 3.0, ready)

	assert.ErrorIs(t, svc.ConsumeFromReady(created.Id, 100), customerrors.ErrInsufficientReady)
}

func TestRecipeService_WasteAndIncrease(t *testing.T) {
	env := newServiceEnv(t)
	svc := newRecipeService(env)

	created, err := svc.InsertNew(models.Product{Name: "Pizza", Price: 12, Ready: 5})
	require.NoError(t, err)

	require.NoError(t, svc.Waste(created.Id, 2, "o-1", "expired", true, models.OrderItem{Id: "i-1"}, "user"))
	ready, err := svc.GetReadyNumber(created.Id)
	require.NoError(t, err)
	assert.Equal(t, 3.0, ready)

	require.NoError(t, svc.Increase(created.Id, 4, "restock", "o-1", "user"))
	ready, err = svc.GetReadyNumber(created.Id)
	require.NoError(t, err)
	assert.Equal(t, 7.0, ready)
}

func TestRecipeService_GetRecipeTree(t *testing.T) {
	env := newServiceEnv(t)
	svc := newRecipeService(env)

	insertMaterial(t, env, models.Material{
		Id:      "mat-tree",
		Name:    "Flour",
		Unit:    "kg",
		Entries: []models.MaterialEntry{{Id: "e1", Quantity: 3}},
	})
	created, err := svc.InsertNew(models.Product{
		Name:      "Pizza",
		Price:     12,
		Materials: []models.Material{{Id: "mat-tree", Quantity: 2}},
	})
	require.NoError(t, err)

	tree, err := svc.GetRecipeTree(created.Id)
	require.NoError(t, err)
	require.Len(t, tree.Materials, 1)
	assert.Equal(t, "Flour", tree.Materials[0].Name)
	assert.Len(t, tree.Materials[0].Entries, 1)
}

func TestRecipeService_CheckRecipesAvailability(t *testing.T) {
	env := newServiceEnv(t)
	svc := newRecipeService(env)

	insertMaterial(t, env, models.Material{
		Id:      "mat-avail",
		Name:    "Flour",
		Unit:    "kg",
		Entries: []models.MaterialEntry{{Id: "e1", Quantity: 10, PurchaseQuantity: 10, PurchasePrice: 10}},
	})
	created, err := svc.InsertNew(models.Product{
		Name:      "Pizza",
		Price:     12,
		Materials: []models.Material{{Id: "mat-avail", Quantity: 2}},
	})
	require.NoError(t, err)

	availabilities, err := svc.CheckRecipesAvailability([]string{created.Id})
	require.NoError(t, err)
	require.Len(t, availabilities, 1)
	assert.Equal(t, created.Id, availabilities[0].RecipeId)
}

func TestRecipeService_FillRecipeDesign(t *testing.T) {
	env := newServiceEnv(t)
	svc := newRecipeService(env)

	created, err := svc.InsertNew(models.Product{Name: "Pizza", Price: 12})
	require.NoError(t, err)

	item, err := svc.FillRecipeDesign(models.OrderItem{Product: models.Product{Id: created.Id}})
	require.NoError(t, err)
	assert.Equal(t, "Pizza", item.Product.Name)
}

func TestRecipeService_GetRecipeMaterials_Error(t *testing.T) {
	env := newServiceEnv(t)
	svc := newRecipeService(env)

	_, err := svc.GetRecipeMaterials("missing")
	assert.Error(t, err)
}
