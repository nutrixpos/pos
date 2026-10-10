package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/nutrixpos/pos/internal/testutil"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func newPOService(env *testutil.TestEnv) *PurchaseOrderService {
	return &PurchaseOrderService{Logger: env.Logger, Config: env.Config, Settings: models.Settings{}}
}

func countDocs(t *testing.T, env *testutil.TestEnv, collection string, filter bson.M) int64 {
	t.Helper()
	n, err := env.Client.Database(env.Config.Databases[0].Database).Collection(collection).CountDocuments(context.Background(), filter)
	require.NoError(t, err)
	return n
}

func makeAutoPO(materialId string, quantity float64) models.PurchaseOrder {
	return models.PurchaseOrder{
		AutoReceive: true,
		Supplier:    "test-supplier",
		Items: []models.PurchaseOrderItem{{
			MaterialId:    materialId,
			Quantity:      quantity,
			PurchasePrice: 2,
		}},
	}
}

func insertMaterialForPO(t *testing.T, env *testutil.TestEnv, id string) {
	t.Helper()
	_, err := env.Client.Database(env.Config.Databases[0].Database).Collection("materials").
		InsertOne(context.Background(), models.Material{Id: id, Name: "Flour", Unit: "kg", Entries: []models.MaterialEntry{}})
	require.NoError(t, err)
}

func TestCreatePurchaseOrder_AutoReceiveSuccess(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterialForPO(t, env, "mat-1")

	svc := newPOService(env)
	po, err := svc.CreatePurchaseOrder(makeAutoPO("mat-1", 5), "user-1")
	require.NoError(t, err)

	assert.Equal(t, models.PurchaseOrderStatusReceived, po.Status)
	require.Len(t, po.Items, 1)
	assert.Equal(t, 5.0, po.Items[0].ReceivedQuantity)
	assert.NotEmpty(t, po.Items[0].EntryIds)

	// Material got one entry of 5 units.
	var material models.Material
	err = env.Client.Database(env.Config.Databases[0].Database).Collection("materials").
		FindOne(context.Background(), bson.M{"id": "mat-1"}).Decode(&material)
	require.NoError(t, err)
	require.Len(t, material.Entries, 1)
	assert.Equal(t, 5.0, material.Entries[0].Quantity)
	// The entry stores the total paid for the batch (unit price * received qty),
	// which is what the cost engine divides by PurchaseQuantity.
	assert.Equal(t, 10.0, material.Entries[0].PurchasePrice)

	assert.Equal(t, int64(1), countDocs(t, env, "grns", bson.M{"purchase_order_id": po.Id}))
	assert.Equal(t, int64(1), countDocs(t, env, "logs", bson.M{"type": models.LogTypeMaterialGRNReceive}))
}

func TestCreatePurchaseOrder_AutoReceiveFailureRollsBack(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterialForPO(t, env, "mat-2")

	testHookAfterPOUpdate = func() error { return fmt.Errorf("injected failure") }
	t.Cleanup(func() { testHookAfterPOUpdate = nil })

	svc := newPOService(env)
	po, err := svc.CreatePurchaseOrder(makeAutoPO("mat-2", 5), "user-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "injected failure")

	// The purchase order is deleted by the caller.
	err = env.Client.Database(env.Config.Databases[0].Database).Collection("purchase_orders").
		FindOne(context.Background(), bson.M{"id": po.Id}).Decode(&models.PurchaseOrder{})
	assert.ErrorIs(t, err, mongo.ErrNoDocuments)

	// No partial receipt data: material has no entries, no logs, no GRNs.
	var material models.Material
	err = env.Client.Database(env.Config.Databases[0].Database).Collection("materials").
		FindOne(context.Background(), bson.M{"id": "mat-2"}).Decode(&material)
	require.NoError(t, err)
	assert.Empty(t, material.Entries)

	assert.Equal(t, int64(0), countDocs(t, env, "logs", bson.M{"type": models.LogTypeMaterialGRNReceive}))
	assert.Equal(t, int64(0), countDocs(t, env, "grns", bson.M{}))
}

func TestReceivePurchaseOrder_FailureRestoresPO(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterialForPO(t, env, "mat-3")

	svc := newPOService(env)
	po, err := svc.CreatePurchaseOrder(models.PurchaseOrder{
		Supplier: "test-supplier",
		Items: []models.PurchaseOrderItem{{
			MaterialId:    "mat-3",
			Quantity:      5,
			PurchasePrice: 2,
		}},
	}, "user-1")
	require.NoError(t, err)
	assert.Equal(t, models.PurchaseOrderStatusOpen, po.Status)

	testHookAfterPOUpdate = func() error { return fmt.Errorf("injected failure") }
	t.Cleanup(func() { testHookAfterPOUpdate = nil })

	_, err = svc.ReceivePurchaseOrder(po.Id, "user-1", nil)
	require.Error(t, err)

	// The purchase order is restored to its pre-receipt state, not deleted.
	var stored models.PurchaseOrder
	err = env.Client.Database(env.Config.Databases[0].Database).Collection("purchase_orders").
		FindOne(context.Background(), bson.M{"id": po.Id}).Decode(&stored)
	require.NoError(t, err)
	assert.Equal(t, models.PurchaseOrderStatusOpen, stored.Status)
	assert.Equal(t, 0.0, stored.Items[0].ReceivedQuantity)
	assert.Empty(t, stored.Items[0].EntryIds)

	var material models.Material
	err = env.Client.Database(env.Config.Databases[0].Database).Collection("materials").
		FindOne(context.Background(), bson.M{"id": "mat-3"}).Decode(&material)
	require.NoError(t, err)
	assert.Empty(t, material.Entries)

	assert.Equal(t, int64(0), countDocs(t, env, "logs", bson.M{"type": models.LogTypeMaterialGRNReceive}))
	assert.Equal(t, int64(0), countDocs(t, env, "grns", bson.M{}))
}

func TestReceivePurchaseOrder_RollbackRemovesPushedEntries(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterialForPO(t, env, "mat-push")
	insertMaterialForPO(t, env, "mat-push-2")

	svc := newPOService(env)
	po, err := svc.CreatePurchaseOrder(models.PurchaseOrder{
		Supplier: "test-supplier",
		Items: []models.PurchaseOrderItem{
			{MaterialId: "mat-push", Quantity: 5, PurchasePrice: 2},
			{MaterialId: "mat-push-2", Quantity: 3, PurchasePrice: 4},
		},
	}, "user-1")
	require.NoError(t, err)
	assert.Equal(t, models.PurchaseOrderStatusOpen, po.Status)

	testHookAfterMaterialPush = func() error { return fmt.Errorf("injected failure after push") }
	t.Cleanup(func() { testHookAfterMaterialPush = nil })

	_, err = svc.ReceivePurchaseOrder(po.Id, "user-1", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "injected failure after push")

	// The entries already pushed to each material were rolled back.
	assert.Empty(t, getMaterial(t, env, "mat-push").Entries)
	assert.Empty(t, getMaterial(t, env, "mat-push-2").Entries)

	// No logs or GRNs survived, and the purchase order was restored.
	assert.Equal(t, int64(0), countDocs(t, env, "logs", bson.M{"type": models.LogTypeMaterialGRNReceive}))
	assert.Equal(t, int64(0), countDocs(t, env, "grns", bson.M{}))

	var stored models.PurchaseOrder
	err = env.Client.Database(env.Config.Databases[0].Database).Collection("purchase_orders").
		FindOne(context.Background(), bson.M{"id": po.Id}).Decode(&stored)
	require.NoError(t, err)
	assert.Equal(t, models.PurchaseOrderStatusOpen, stored.Status)
	assert.Equal(t, 0.0, stored.Items[0].ReceivedQuantity)
	assert.Empty(t, stored.Items[0].EntryIds)
}

func TestReceivePurchaseOrder_PartialThenFull(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterialForPO(t, env, "mat-a")
	insertMaterialForPO(t, env, "mat-b")

	svc := newPOService(env)
	po, err := svc.CreatePurchaseOrder(models.PurchaseOrder{
		Supplier: "test-supplier",
		Items: []models.PurchaseOrderItem{
			{MaterialId: "mat-a", Quantity: 10, PurchasePrice: 1},
			{MaterialId: "mat-b", Quantity: 10, PurchasePrice: 1},
		},
	}, "user-1")
	require.NoError(t, err)

	// Partial receive of one item.
	grn, err := svc.ReceivePurchaseOrder(po.Id, "user-1", []ReceiveItem{{ItemId: po.Items[0].ItemId, Quantity: 4}})
	require.NoError(t, err)
	assert.NotEmpty(t, grn.Id)

	var stored models.PurchaseOrder
	err = env.Client.Database(env.Config.Databases[0].Database).Collection("purchase_orders").
		FindOne(context.Background(), bson.M{"id": po.Id}).Decode(&stored)
	require.NoError(t, err)
	assert.Equal(t, models.PurchaseOrderStatusPartial, stored.Status)
	assert.Equal(t, 4.0, stored.Items[0].ReceivedQuantity)

	// Full receive of the remainder.
	_, err = svc.ReceivePurchaseOrder(po.Id, "user-1", nil)
	require.NoError(t, err)

	err = env.Client.Database(env.Config.Databases[0].Database).Collection("purchase_orders").
		FindOne(context.Background(), bson.M{"id": po.Id}).Decode(&stored)
	require.NoError(t, err)
	assert.Equal(t, models.PurchaseOrderStatusReceived, stored.Status)
	assert.Equal(t, 10.0, stored.Items[0].ReceivedQuantity)
	assert.Equal(t, 10.0, stored.Items[1].ReceivedQuantity)

	var materialA models.Material
	err = env.Client.Database(env.Config.Databases[0].Database).Collection("materials").
		FindOne(context.Background(), bson.M{"id": "mat-a"}).Decode(&materialA)
	require.NoError(t, err)
	var sum float64
	for _, e := range materialA.Entries {
		sum += e.Quantity
		// Each received entry must keep the PO unit price (1) as its per-unit cost,
		// regardless of how the receipt was split.
		assert.InDelta(t, 1.0, e.PurchasePrice/e.PurchaseQuantity, 0.0001)
	}
	assert.Equal(t, 10.0, sum)
}

func TestPurchaseOrder_ValidationErrors(t *testing.T) {
	env := newServiceEnv(t)
	svc := newPOService(env)

	_, err := svc.CreatePurchaseOrder(models.PurchaseOrder{Items: []models.PurchaseOrderItem{{MaterialId: "", Quantity: 1}}}, "u")
	assert.Error(t, err)
	_, err = svc.CreatePurchaseOrder(models.PurchaseOrder{Items: []models.PurchaseOrderItem{{MaterialId: "x", Quantity: 0}}}, "u")
	assert.Error(t, err)
	_, err = svc.CreatePurchaseOrder(models.PurchaseOrder{Items: []models.PurchaseOrderItem{{MaterialId: "x", Quantity: 1, PurchasePrice: -1}}}, "u")
	assert.Error(t, err)
	_, err = svc.CreatePurchaseOrder(models.PurchaseOrder{Items: []models.PurchaseOrderItem{{MaterialId: "nope", Quantity: 1}}}, "u")
	assert.Error(t, err)

	_, err = svc.ReceivePurchaseOrder("missing", "u", nil)
	assert.Error(t, err)
}

func TestCancelPurchaseOrder_ReceivedError(t *testing.T) {
	env := newServiceEnv(t)
	insertMaterialForPO(t, env, "mat-cancel-rec")
	svc := newPOService(env)
	po, err := svc.CreatePurchaseOrder(makeAutoPO("mat-cancel-rec", 1), "u")
	require.NoError(t, err)
	assert.Error(t, svc.CancelPurchaseOrder(po.Id))
}

func TestReceivePurchaseOrder_AlreadyReceived(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterialForPO(t, env, "mat-4")

	svc := newPOService(env)
	po, err := svc.CreatePurchaseOrder(makeAutoPO("mat-4", 5), "user-1")
	require.NoError(t, err)
	assert.Equal(t, models.PurchaseOrderStatusReceived, po.Status)

	_, err = svc.ReceivePurchaseOrder(po.Id, "user-1", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already fully received")
}

// TestReceivedPOEntryProducesCorrectOrderCost guards the seam between purchase
// order receipt and the order cost engine: a received PO must store the batch
// total so the derived unit cost equals the PO unit price, otherwise order cost
// and the profit shown in sales reports are wrong.
func TestReceivedPOEntryProducesCorrectOrderCost(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterialForPO(t, env, "mat-cost")

	posvc := newPOService(env)
	_, err := posvc.CreatePurchaseOrder(models.PurchaseOrder{
		AutoReceive: true,
		Supplier:    "test-supplier",
		Items: []models.PurchaseOrderItem{{
			MaterialId:    "mat-cost",
			Quantity:      5,
			PurchasePrice: 2,
		}},
	}, "user-1")
	require.NoError(t, err)

	var material models.Material
	err = env.Client.Database(env.Config.Databases[0].Database).Collection("materials").
		FindOne(context.Background(), bson.M{"id": "mat-cost"}).Decode(&material)
	require.NoError(t, err)
	require.Len(t, material.Entries, 1)

	_, err = env.Client.Database(env.Config.Databases[0].Database).Collection("recipes").
		InsertOne(context.Background(), models.Product{Id: "prod-cost", Name: "Pizza", Price: 20})
	require.NoError(t, err)

	orderItem := models.OrderItem{
		Id:       "item-cost",
		Product:  models.Product{Id: "prod-cost"},
		Quantity: 2,
		Materials: []models.OrderItemMaterial{{
			Material: material,
			Entry:    material.Entries[0],
			Quantity: 3,
		}},
	}

	svc := &OrderService{Logger: env.Logger, Config: env.Config}
	costs, err := svc.CalculateCost([]models.OrderItem{orderItem})
	require.NoError(t, err)
	require.Len(t, costs, 1)

	// unit price 2 * component qty 3 * item qty 2 = 12.
	assert.Equal(t, 12.0, costs[0].Cost)
	assert.Equal(t, 40.0, costs[0].SalePrice)
}
