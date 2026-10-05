package services

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/nutrixpos/pos/internal/testutil"
	"github.com/nutrixpos/pos/modules/core/dto"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGetOrderDisplayId_ConcurrentUnique(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())

	const workers = 25
	type result struct {
		id  string
		err error
	}
	results := make(chan result, workers)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			svc := &OrderService{Logger: env.Logger, Config: env.Config}
			id, err := svc.GetOrderDisplayId()
			results <- result{id: id, err: err}
		}()
	}
	wg.Wait()
	close(results)

	seen := map[string]bool{}
	for r := range results {
		require.NoError(t, r.err)
		require.NotEmpty(t, r.id)
		require.False(t, seen[r.id], "duplicate display id %s", r.id)
		seen[r.id] = true
	}
	require.Len(t, seen, workers)
}

// TestRefundItem_RetrySkipsInventoryReturn ensures a retried refund does not
// return the same InventoryReturnQty to stock twice after the first attempt
// already completed the inventory-return step.
func TestRefundItem_RetrySkipsInventoryReturn(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())

	material := models.Material{
		Id:   "mat-refund",
		Name: "Flour",
		Entries: []models.MaterialEntry{{
			Id: "entry-1", Quantity: 5, PurchasePrice: 10, PurchaseQuantity: 10,
		}},
	}
	_, err := env.Client.Database(env.Config.Databases[0].Database).Collection("materials").
		InsertOne(context.Background(), material)
	require.NoError(t, err)

	order := models.Order{
		Id: "order-1",
		Items: []models.OrderItem{{
			Id: "item-1",
			Materials: []models.OrderItemMaterial{{
				Material: material,
				Entry:    material.Entries[0],
				Quantity: 2,
			}},
		}},
	}
	_, err = env.Client.Database(env.Config.Databases[0].Database).Collection("orders").
		InsertOne(context.Background(), order)
	require.NoError(t, err)

	request := dto.OrderItemRefundRequest{
		OrderId:     "order-1",
		ItemId:      "item-1",
		Reason:      "damaged",
		RefundValue: 5,
		Destination: dto.DTOOrderItemRefundDestination_Custom,
		MaterialRefunds: []dto.OrderItemRefundMaterialDTO{{
			MaterialId:         "mat-refund",
			EntryId:            "entry-1",
			InventoryReturnQty: 2,
			DisposeQty:         1,
		}},
	}

	svc := &OrderService{Logger: env.Logger, Config: env.Config, Settings: models.Settings{}}

	require.NoError(t, svc.RefundItem(request, "user-1"))

	entryQuantity := func() float64 {
		var m models.Material
		err := env.Client.Database(env.Config.Databases[0].Database).Collection("materials").
			FindOne(context.Background(), bson.M{"id": "mat-refund"}).Decode(&m)
		require.NoError(t, err)
		require.Len(t, m.Entries, 1)
		return m.Entries[0].Quantity
	}

	disposalCount := func() int64 {
		n, err := env.Client.Database(env.Config.Databases[0].Database).Collection("disposals").
			CountDocuments(context.Background(), bson.M{"disposal.order_id": "order-1", "disposal.type": models.TypeDisposalMaterial})
		require.NoError(t, err)
		return n
	}

	salesDay := func() models.SalesPerDay {
		var s models.SalesPerDay
		err := env.Client.Database(env.Config.Databases[0].Database).Collection(env.Config.Databases[0].Tables["sales"]).
			FindOne(context.Background(), bson.M{"date": time.Now().Format("2006-01-02")}).Decode(&s)
		require.NoError(t, err)
		return s
	}

	assert.Equal(t, 7.0, entryQuantity()) // 5 + 2 returned once
	assert.Equal(t, int64(1), disposalCount())

	day := salesDay()
	assert.Len(t, day.Refunds, 1)
	assert.Equal(t, 5.0, day.RefundsValue)

	// Simulate a client retry after a completed refund: no effect may repeat.
	require.NoError(t, svc.RefundItem(request, "user-1"))
	assert.Equal(t, 7.0, entryQuantity())
	assert.Equal(t, int64(1), disposalCount())

	day = salesDay()
	assert.Len(t, day.Refunds, 1)
	assert.Equal(t, 5.0, day.RefundsValue)
}
