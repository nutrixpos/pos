package services

import (
	"testing"

	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaterialService_Availability(t *testing.T) {
	env := newServiceEnv(t)
	insertMaterial(t, env, models.Material{
		Id: "mat-avail2",
		Entries: []models.MaterialEntry{
			{Id: "e1", Quantity: 4},
			{Id: "e2", Quantity: 6},
		},
	})
	svc := newMaterialService(env, models.Settings{})

	total, err := svc.GetComponentAvailability("mat-avail2")
	require.NoError(t, err)
	assert.Equal(t, 10.0, total)

	entry, err := svc.GetMaterialEntryAvailability("mat-avail2", "e1")
	require.NoError(t, err)
	assert.Equal(t, 4.0, entry)

	_, err = svc.GetMaterialEntryAvailability("mat-avail2", "missing")
	assert.Error(t, err)
}

func TestMaterialService_ConsumeFromInventory(t *testing.T) {
	env := newServiceEnv(t)
	insertMaterial(t, env, models.Material{
		Id:      "mat-consume-inv",
		Name:    "Flour",
		Entries: []models.MaterialEntry{{Id: "e1", Quantity: 50}},
	})
	svc := newMaterialService(env, models.Settings{Inventory: models.MaterialSettings{StockAlertTreshold: 1000}})

	notifications, err := svc.ConsumeFromInventory(models.Material{Id: "mat-consume-inv", Name: "Flour"}, "e1", 5, "order", "o-1", "u")
	require.NoError(t, err)
	assert.NotEmpty(t, notifications) // low stock notification (threshold high)

	// Insufficient quantity is rejected.
	_, err = svc.ConsumeFromInventory(models.Material{Id: "mat-consume-inv", Name: "Flour"}, "e1", 1000, "order", "o-1", "u")
	assert.Error(t, err)
}

func TestRecipeService_SubProducts(t *testing.T) {
	env := newServiceEnv(t)
	svc := newRecipeService(env)

	sub, err := svc.InsertNew(models.Product{Name: "Sauce", Price: 2})
	require.NoError(t, err)
	parent, err := svc.InsertNew(models.Product{
		Name:        "Pizza",
		Price:       12,
		SubProducts: []models.Product{{Id: sub.Id, Name: "Sauce", Quantity: 1}},
	})
	require.NoError(t, err)

	products, _, err := svc.GetProducts(GetProductsParams{PageNumber: 1, PageSize: 10})
	require.NoError(t, err)
	assert.NotEmpty(t, products)

	tree, err := svc.GetRecipeTree(parent.Id)
	require.NoError(t, err)
	require.Len(t, tree.SubProducts, 1)
	assert.Equal(t, sub.Id, tree.SubProducts[0].Id)

	availabilities, err := svc.CheckRecipesAvailability([]string{parent.Id})
	require.NoError(t, err)
	require.Len(t, availabilities, 1)
}

func TestOrderService_CalculateCost_Branches(t *testing.T) {
	env := newServiceEnv(t)
	insertMaterial(t, env, models.Material{
		Id:      "m-branch",
		Name:    "Flour",
		Entries: []models.MaterialEntry{{Id: "e1", Quantity: 10, PurchaseQuantity: 10, PurchasePrice: 100}},
	})
	recipeSvc := &RecipeService{Logger: env.Logger, Config: env.Config}
	created, err := recipeSvc.InsertNew(models.Product{Name: "Pizza", Price: 10, Materials: []models.Material{{Id: "m-branch", Quantity: 1}}})
	require.NoError(t, err)

	// Entry referenced by the order item does not exist on the material.
	svc := newOrderSvc(env, "exact")
	_, err = svc.CalculateCost([]models.OrderItem{{
		Id:        "i-1",
		Product:   models.Product{Id: created.Id},
		Quantity:  1,
		Materials: []models.OrderItemMaterial{{Material: models.Material{Id: "m-branch"}, Entry: models.MaterialEntry{Id: "nope"}, Quantity: 1}},
	}})
	assert.Error(t, err)

	// Average with a material that has no entries: no panic.
	insertMaterial(t, env, models.Material{Id: "m-empty", Name: "Empty"})
	emptyRecipe, err := recipeSvc.InsertNew(models.Product{Name: "Empty", Price: 1, Materials: []models.Material{{Id: "m-empty", Quantity: 1}}})
	require.NoError(t, err)
	_, err = newOrderSvc(env, "average").CalculateCost([]models.OrderItem{{
		Id:        "i-2",
		Product:   models.Product{Id: emptyRecipe.Id},
		Quantity:  1,
		Materials: []models.OrderItemMaterial{{Material: models.Material{Id: "m-empty"}, Entry: models.MaterialEntry{Id: "e1"}, Quantity: 1}},
	}})
	require.NoError(t, err)
}

func TestOrderService_PrintReceipt_Error(t *testing.T) {
	env := newServiceEnv(t)
	svc := newOrderSvc(env, "average")
	// No printer reachable on the configured host.
	assert.Error(t, svc.PrintReceipt(models.Order{Id: "o-1"}, "template", "en", "127.0.0.1:1"))
}

func TestOrderService_ConsumeOrderComponents_NoConsumption(t *testing.T) {
	env := newServiceEnv(t)
	svc := newOrderSvc(env, "average")
	// Item without inventory consumption short-circuits the loop.
	order := models.Order{Id: "o-nc", Items: []models.OrderItem{{
		Id:      "i-1",
		Product: models.Product{EnableInventoryConsumption: false},
	}}}
	require.NoError(t, svc.ConsumeOrderComponents(order, "u"))
}

func TestRecipeService_FillRecipeDesign_SubItems(t *testing.T) {
	env := newServiceEnv(t)
	svc := newRecipeService(env)

	sub, err := svc.InsertNew(models.Product{Name: "Sauce", Price: 2})
	require.NoError(t, err)
	parent, err := svc.InsertNew(models.Product{
		Name:        "Pizza",
		Price:       12,
		SubProducts: []models.Product{{Id: sub.Id, Quantity: 1}},
	})
	require.NoError(t, err)

	item, err := svc.FillRecipeDesign(models.OrderItem{
		Product:  models.Product{Id: parent.Id},
		SubItems: []models.OrderItem{{Product: models.Product{Id: sub.Id}, Quantity: 1}},
	})
	require.NoError(t, err)
	assert.Equal(t, "Pizza", item.Product.Name)
	require.Len(t, item.SubItems, 1)
	assert.Equal(t, "Sauce", item.SubItems[0].Product.Name)
}

func TestMaterialService_ConsumeFromReadyItem(t *testing.T) {
	env := newServiceEnv(t)
	recipeSvc := newRecipeService(env)
	created, err := recipeSvc.InsertNew(models.Product{Name: "Pizza", Price: 10, Ready: 5})
	require.NoError(t, err)

	svc := newMaterialService(env, models.Settings{})
	_, err = svc.ConsumeItemComponentsForOrder(
		models.OrderItem{Product: models.Product{Id: created.Id}, Quantity: 2, IsConsumeFromReady: true},
		models.Order{Id: "o"}, 0, "u")
	require.NoError(t, err)

	ready, err := recipeSvc.GetReadyNumber(created.Id)
	require.NoError(t, err)
	assert.Equal(t, 3.0, ready)
}

func TestRecipeService_GetProducts_DanglingSubProduct(t *testing.T) {
	env := newServiceEnv(t)
	svc := newRecipeService(env)
	_, err := svc.InsertNew(models.Product{Name: "Pizza", Price: 10, SubProducts: []models.Product{{Id: "missing"}}})
	require.NoError(t, err)

	_, _, err = svc.GetProducts(GetProductsParams{PageNumber: 1, PageSize: 10})
	assert.Error(t, err)
}

func TestSeeder_NoMaterials(t *testing.T) {
	env := newServiceEnv(t)
	svc := Seeder{Logger: env.Logger, Config: env.Config, Prompter: fakePrompter{confirm: false}}

	require.NoError(t, svc.SeedProducts())
	require.NoError(t, svc.SeedCategories())
}
