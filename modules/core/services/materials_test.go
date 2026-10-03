package services

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nutrixpos/pos/internal/testutil"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func materialsCollection(t *testing.T, env *testutil.TestEnv) *mongo.Collection {
	t.Helper()
	return env.Client.Database(env.Config.Databases[0].Database).Collection("materials")
}

func insertMaterial(t *testing.T, env *testutil.TestEnv, material models.Material) {
	t.Helper()
	_, err := materialsCollection(t, env).InsertOne(context.Background(), material)
	require.NoError(t, err)
}

func getMaterial(t *testing.T, env *testutil.TestEnv, id string) models.Material {
	t.Helper()
	var m models.Material
	err := materialsCollection(t, env).FindOne(context.Background(), bson.M{"id": id}).Decode(&m)
	require.NoError(t, err)
	return m
}

func newMaterialService(env *testutil.TestEnv, settings models.Settings) *MaterialService {
	return &MaterialService{Logger: env.Logger, Config: env.Config, Settings: settings}
}

func TestGetMaterialEntries_Pagination(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())

	entries := make([]models.MaterialEntry, 0, 7)
	for i := 0; i < 7; i++ {
		entries = append(entries, models.MaterialEntry{Id: fmt.Sprintf("entry-%d", i), Quantity: float64(i)})
	}
	insertMaterial(t, env, models.Material{Id: "mat-pagination", Entries: entries})

	svc := newMaterialService(env, models.Settings{})

	tests := []struct {
		name       string
		pageNumber int
		pageSize   int
		want       int
	}{
		{"first page", 0, 2, 2},
		{"second page", 1, 2, 2},
		{"last page partial", 3, 2, 1},
		{"beyond end", 10, 2, 0},
		{"zero page size", 0, 0, 0},
		{"full page", 0, 7, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, total, err := svc.GetMaterialEntries("mat-pagination", GetMaterialEntriesParams{PageNumber: tt.pageNumber, PageSize: tt.pageSize})
			require.NoError(t, err)
			assert.Len(t, got, tt.want)
			assert.Equal(t, int64(7), total)
		})
	}
}

func TestGetMaterialEntries_RejectsNegativeParams(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterial(t, env, models.Material{Id: "mat-neg", Entries: []models.MaterialEntry{{Id: "e", Quantity: 1}}})

	svc := newMaterialService(env, models.Settings{})

	_, _, err := svc.GetMaterialEntries("mat-neg", GetMaterialEntriesParams{PageNumber: -1, PageSize: 2})
	require.Error(t, err)

	_, _, err = svc.GetMaterialEntries("mat-neg", GetMaterialEntriesParams{PageNumber: 0, PageSize: -2})
	require.Error(t, err)
}

func TestWaste_NonConsume(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterial(t, env, models.Material{Id: "mat-waste", Entries: []models.MaterialEntry{{Id: "entry-1", Quantity: 10}}})

	svc := newMaterialService(env, models.Settings{})
	require.NoError(t, svc.Waste("entry-1", "mat-waste", 3, "order-1", "expired", false, "user-1"))

	got := getMaterial(t, env, "mat-waste")
	require.Len(t, got.Entries, 1)
	assert.Equal(t, 7.0, got.Entries[0].Quantity)
}

func TestWaste_Consume(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterial(t, env, models.Material{Id: "mat-consume", Entries: []models.MaterialEntry{{Id: "entry-1", Quantity: 10}}})

	svc := newMaterialService(env, models.Settings{})
	require.NoError(t, svc.Waste("entry-1", "mat-consume", 3, "order-1", "cooked", true, "user-1"))

	got := getMaterial(t, env, "mat-consume")
	require.Len(t, got.Entries, 1)
	assert.Equal(t, 7.0, got.Entries[0].Quantity)
}

func TestWaste_UnknownEntry(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterial(t, env, models.Material{Id: "mat-missing", Entries: []models.MaterialEntry{{Id: "entry-1", Quantity: 10}}})

	svc := newMaterialService(env, models.Settings{})
	err := svc.Waste("entry-404", "mat-missing", 1, "order-1", "x", false, "user-1")
	require.Error(t, err)
}

func TestInventoryReturn(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterial(t, env, models.Material{Id: "mat-ret", Entries: []models.MaterialEntry{{Id: "entry-1", Quantity: 5}}})

	svc := newMaterialService(env, models.Settings{})
	require.NoError(t, svc.InventoryReturn("entry-1", "mat-ret", 4, "order-1", "return", false, "user-1"))

	got := getMaterial(t, env, "mat-ret")
	require.Len(t, got.Entries, 1)
	assert.Equal(t, 9.0, got.Entries[0].Quantity)
}

func TestConsumeItemComponentsForOrder_Exact(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	material := models.Material{
		Id:   "mat-exact",
		Name: "Flour",
		Entries: []models.MaterialEntry{{
			Id: "entry-1", Quantity: 10, PurchasePrice: 100, PurchaseQuantity: 10,
		}},
	}
	insertMaterial(t, env, material)

	svc := newMaterialService(env, models.Settings{Orders: models.OrderSettings{DefaultCostCalculationMethod: "exact"}})

	item := models.OrderItem{
		Product:  models.Product{Id: "prod-1"},
		Quantity: 2,
		Materials: []models.OrderItemMaterial{{
			Material: material,
			Entry:    material.Entries[0],
			Quantity: 2,
		}},
	}
	order := models.Order{Id: "order-1", DisplayId: "A-1"}

	_, err := svc.ConsumeItemComponentsForOrder(item, order, 0, "user-1")
	require.NoError(t, err)

	got := getMaterial(t, env, "mat-exact")
	require.Len(t, got.Entries, 1)
	assert.Equal(t, 6.0, got.Entries[0].Quantity) // 10 - (2 * 2)
}

func TestConsumeItemComponentsForOrder_ExactInsufficient(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	material := models.Material{
		Id:   "mat-insuff",
		Name: "Flour",
		Entries: []models.MaterialEntry{{
			Id: "entry-1", Quantity: 1, PurchasePrice: 100, PurchaseQuantity: 10,
		}},
	}
	insertMaterial(t, env, material)

	svc := newMaterialService(env, models.Settings{Orders: models.OrderSettings{DefaultCostCalculationMethod: "exact"}})

	item := models.OrderItem{
		Product:  models.Product{Id: "prod-1"},
		Quantity: 1,
		Materials: []models.OrderItemMaterial{{
			Material: material,
			Entry:    material.Entries[0],
			Quantity: 3,
		}},
	}

	notifications, err := svc.ConsumeItemComponentsForOrder(item, models.Order{Id: "order-1"}, 0, "user-1")
	require.Error(t, err)
	assert.Len(t, notifications, 1)
	assert.Equal(t, "inventory_insufficient", notifications[0].TopicName)

	// Nothing was consumed.
	got := getMaterial(t, env, "mat-insuff")
	assert.Equal(t, 1.0, got.Entries[0].Quantity)
}

func TestConsumeItemComponentsForOrder_Average(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	now := time.Now()
	material := models.Material{
		Id:   "mat-avg",
		Name: "Flour",
		Entries: []models.MaterialEntry{
			{Id: "entry-1", Quantity: 3, ExpirationDate: now.Add(72 * time.Hour)},
			{Id: "entry-2", Quantity: 5, ExpirationDate: now.Add(24 * time.Hour)},
			{Id: "entry-3", Quantity: 10, ExpirationDate: now.Add(48 * time.Hour)},
		},
	}
	insertMaterial(t, env, material)

	svc := newMaterialService(env, models.Settings{Orders: models.OrderSettings{DefaultCostCalculationMethod: "average"}})

	item := models.OrderItem{
		Product:  models.Product{Id: "prod-1"},
		Quantity: 1,
		Materials: []models.OrderItemMaterial{{
			Material: material,
			Entry:    material.Entries[0],
			Quantity: 6,
		}},
	}

	_, err := svc.ConsumeItemComponentsForOrder(item, models.Order{Id: "order-1", DisplayId: "A-1"}, 0, "user-1")
	require.NoError(t, err)

	got := getMaterial(t, env, "mat-avg")
	quantities := map[string]float64{}
	for _, e := range got.Entries {
		quantities[e.Id] = e.Quantity
	}
	// FIFO by expiration: entry-2 (5) then entry-3 (1); entry-1 untouched.
	assert.Equal(t, 3.0, quantities["entry-1"])
	assert.Equal(t, 0.0, quantities["entry-2"])
	assert.Equal(t, 9.0, quantities["entry-3"])
}

func TestConcurrentWaste_NoLostUpdates(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterial(t, env, models.Material{Id: "mat-conc", Entries: []models.MaterialEntry{{Id: "entry-1", Quantity: 100}}})

	const workers = 10
	errs := make(chan error, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			svc := newMaterialService(env, models.Settings{})
			errs <- svc.Waste("entry-1", "mat-conc", 1, "order-1", "x", false, "user-1")
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		require.NoError(t, err)
	}

	got := getMaterial(t, env, "mat-conc")
	require.Len(t, got.Entries, 1)
	assert.Equal(t, 90.0, got.Entries[0].Quantity)
}
