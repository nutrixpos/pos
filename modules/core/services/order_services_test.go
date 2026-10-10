package services

import (
	"context"
	"testing"
	"time"

	"github.com/nutrixpos/pos/common/customerrors"
	"github.com/nutrixpos/pos/internal/testutil"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func newOrderSvc(env *testutil.TestEnv, method string) *OrderService {
	return &OrderService{
		Logger: env.Logger,
		Config: env.Config,
		Settings: models.Settings{
			Inventory: models.MaterialSettings{StockAlertTreshold: 1000},
			Orders: models.OrderSettings{
				DefaultCostCalculationMethod: method,
				Queues:                       []models.OrderQueueSettings{{Prefix: "A", Next: 1}},
			},
		},
	}
}

func TestOrderService_Tips(t *testing.T) {
	env := newServiceEnv(t)
	seedServiceDoc(t, env, "orders", models.Order{Id: "o-tip", Tips: 0})
	svc := newOrderSvc(env, "average")

	require.NoError(t, svc.AddTip("o-tip", 5))
	assert.Equal(t, 5.0, getOrderTips(t, env, "o-tip"))

	require.NoError(t, svc.RemoveTip("o-tip", 2))
	assert.Equal(t, 3.0, getOrderTips(t, env, "o-tip"))
}

func TestOrderService_SimpleMutations(t *testing.T) {
	env := newServiceEnv(t)
	seedServiceDoc(t, env, "orders", models.Order{Id: "o-1", State: "pending"})
	seedServiceDoc(t, env, "orders", models.Order{Id: "o-2", IsPayLater: true, IsPaid: false, State: "pending"})
	svc := newOrderSvc(env, "average")

	require.NoError(t, svc.WasteOrderItem(models.OrderItem{Id: "i-1"}, "o-1", 1, "burnt", map[string]interface{}{}, "u"))
	require.NoError(t, svc.StartOrder("o-1", []models.OrderItem{{Id: "i-1"}}, "u"))
	require.NoError(t, svc.UpdateCustomData("o-1", map[string]string{"table": "5"}))

	var started models.Order
	require.NoError(t, env.Client.Database(env.Config.Databases[0].Database).Collection("orders").FindOne(context.Background(), bson.M{"id": "o-1"}).Decode(&started))
	assert.Equal(t, "in_progress", started.State)
	assert.Equal(t, "5", started.CustomData["table"])

	unpaid, err := svc.GetUnpaidOrders()
	require.NoError(t, err)
	assert.Len(t, unpaid, 1)

	require.NoError(t, svc.CancelOrder("o-2"))
	var cancelled models.Order
	require.NoError(t, env.Client.Database(env.Config.Databases[0].Database).Collection("orders").FindOne(context.Background(), bson.M{"id": "o-2"}).Decode(&cancelled))
	assert.Equal(t, "cancelled", cancelled.State)

	require.NoError(t, svc.DeleteOrder("o-2"))
	_, err = svc.GetOrder("o-2")
	assert.Error(t, err)
}

func TestOrderService_GetLogs(t *testing.T) {
	env := newServiceEnv(t)
	seedServiceDoc(t, env, "logs", map[string]interface{}{"id": "l-1", "order_id": "o-logs", "type": models.LogTypeOrderStart})
	svc := newOrderSvc(env, "average")

	logs, err := svc.GetLogs("o-logs")
	require.NoError(t, err)
	assert.Len(t, logs, 1)
}

func TestOrderService_PayUnpaidOrder(t *testing.T) {
	env := newServiceEnv(t)
	seedServiceDoc(t, env, "orders", models.Order{Id: "o-pay", SalePrice: 10, IsPaid: false})
	svc := newOrderSvc(env, "average")

	assert.Error(t, svc.PayUnpaidOrder("o-pay", nil))
	assert.Error(t, svc.PayUnpaidOrder("o-pay", []models.OrderPayment{{Source: "", Amount: 10}}))
	assert.Error(t, svc.PayUnpaidOrder("o-pay", []models.OrderPayment{{Source: "Cash", Amount: 0}}))
	assert.Error(t, svc.PayUnpaidOrder("o-pay", []models.OrderPayment{{Source: "Cash", Amount: 5}}))
	require.NoError(t, svc.PayUnpaidOrder("o-pay", []models.OrderPayment{{Source: "Cash", Amount: 10}}))

	// Paying twice returns ErrOrderAlreadyPaid.
	err := svc.PayUnpaidOrder("o-pay", []models.OrderPayment{{Source: "Cash", Amount: 10}})
	assert.ErrorIs(t, err, customerrors.ErrOrderAlreadyPaid)
}

func TestOrderService_SubmitAndFinish(t *testing.T) {
	env := newServiceEnv(t)
	recipeSvc := &RecipeService{Logger: env.Logger, Config: env.Config}
	created, err := recipeSvc.InsertNew(models.Product{Name: "Pizza", Price: 10})
	require.NoError(t, err)

	svc := newOrderSvc(env, "average")
	order, err := svc.SubmitOrder(models.Order{
		IsPayLater: true,
		Items:      []models.OrderItem{{Product: models.Product{Id: created.Id}, Quantity: 2}},
	})
	require.NoError(t, err)
	assert.NotEmpty(t, order.Id)
	assert.Equal(t, 20.0, order.SalePrice)

	require.NoError(t, svc.StartOrder(order.Id, order.Items, "u"))
	require.NoError(t, svc.FinishOrder(order.Id, "u"))

	finished, err := svc.GetOrder(order.Id)
	require.NoError(t, err)
	assert.Equal(t, "finished", finished.State)
}

func TestOrderService_GetOrders(t *testing.T) {
	env := newServiceEnv(t)
	seedServiceDoc(t, env, "orders", models.Order{Id: "o-1", DisplayId: "A-1", State: "pending", IsPaid: false})
	seedServiceDoc(t, env, "orders", models.Order{Id: "o-2", DisplayId: "A-2", State: "finished", IsPaid: true})
	svc := newOrderSvc(env, "average")

	orders, total, err := svc.GetOrders(GetOrdersParameters{PageNumber: 1, PageSize: 10, FilterIsPaid: 1, IsPayLater: -1})
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, orders, 1)

	orders, _, err = svc.GetOrders(GetOrdersParameters{PageNumber: 1, PageSize: 10, FilterIsPaid: -1, IsPayLater: 1, OrderDisplayIdContains: "A-1", FilterState: []string{"!finished"}})
	require.NoError(t, err)
	assert.Empty(t, orders)
}

func TestOrderService_CalculateCostMethods(t *testing.T) {
	env := newServiceEnv(t)
	insertMaterial(t, env, models.Material{
		Id:   "m-cost",
		Name: "Flour",
		Entries: []models.MaterialEntry{
			{Id: "e1", Quantity: 10, PurchaseQuantity: 10, PurchasePrice: 100},
			{Id: "e2", Quantity: 10, PurchaseQuantity: 10, PurchasePrice: 200},
		},
	})
	recipeSvc := &RecipeService{Logger: env.Logger, Config: env.Config}
	created, err := recipeSvc.InsertNew(models.Product{Name: "Pizza", Price: 50, Materials: []models.Material{{Id: "m-cost", Quantity: 2}}})
	require.NoError(t, err)

	item := models.OrderItem{
		Id:       "i-1",
		Product:  models.Product{Id: created.Id},
		Quantity: 2,
		Materials: []models.OrderItemMaterial{{
			Material: models.Material{Id: "m-cost", Name: "Flour"},
			Entry:    models.MaterialEntry{Id: "e1"},
			Quantity: 1,
		}},
	}

	// exact: (100/10) * 1 * 2 = 20
	costs, err := newOrderSvc(env, "exact").CalculateCost([]models.OrderItem{item})
	require.NoError(t, err)
	require.Len(t, costs, 1)
	assert.Equal(t, 20.0, costs[0].Cost)

	// average: ((100/10)+(200/10))/2 * 1 * 2 = 30
	costs, err = newOrderSvc(env, "average").CalculateCost([]models.OrderItem{item})
	require.NoError(t, err)
	require.Len(t, costs, 1)
	assert.Equal(t, 30.0, costs[0].Cost)

	// fixed: 3 * 2 = 6
	fixedRecipe, err := recipeSvc.InsertNew(models.Product{Name: "Fixed", Price: 50, EnableFixedCost: true, FixedCost: 3})
	require.NoError(t, err)
	item.Product = models.Product{Id: fixedRecipe.Id}
	costs, err = newOrderSvc(env, "exact").CalculateCost([]models.OrderItem{item})
	require.NoError(t, err)
	require.Len(t, costs, 1)
	assert.Equal(t, 6.0, costs[0].Cost)
	assert.Equal(t, "fixed", costs[0].CostMethod)
}

func TestOrderService_ConsumeOrderComponents(t *testing.T) {
	env := newServiceEnv(t)
	insertMaterial(t, env, models.Material{
		Id:   "m-consume",
		Name: "Flour",
		Entries: []models.MaterialEntry{{
			Id: "e1", Quantity: 10, PurchaseQuantity: 10, PurchasePrice: 10, ExpirationDate: nowTime().Add(time.Hour),
		}},
	})
	svc := newOrderSvc(env, "average")
	order := models.Order{Id: "o-consume"}
	item := models.OrderItem{
		Id:       "i-1",
		Product:  models.Product{Name: "Pizza", EnableInventoryConsumption: true},
		Quantity: 1,
		Materials: []models.OrderItemMaterial{{
			Material: models.Material{Id: "m-consume", Name: "Flour"},
			Entry:    models.MaterialEntry{Id: "e1"},
			Quantity: 1,
		}},
	}
	order.Items = []models.OrderItem{item}
	require.NoError(t, svc.ConsumeOrderComponents(order, "u"))

	got := getMaterial(t, env, "m-consume")
	assert.Less(t, got.Entries[0].Quantity, 10.0)
}

func TestOrderService_MissingEntities(t *testing.T) {
	env := newServiceEnv(t)
	svc := newOrderSvc(env, "average")

	assert.Error(t, svc.AddTip("missing", 1))
	assert.Error(t, svc.RemoveTip("missing", 1))
	assert.Error(t, svc.StartOrder("missing", nil, "u"))
	assert.Error(t, svc.FinishOrder("missing", "u"))
	assert.Error(t, svc.PayUnpaidOrder("missing", []models.OrderPayment{{Source: "Cash", Amount: 1}}))
	_, err := svc.GetOrder("missing")
	assert.Error(t, err)

	// SubmitOrder rejects an unpaid order that carries no payments and is not pay-later.
	_, err = svc.SubmitOrder(models.Order{Items: []models.OrderItem{}})
	assert.Error(t, err)
}

func TestOrderService_GetOrders_CostMethodDefaults(t *testing.T) {
	env := newServiceEnv(t)
	seedServiceDoc(t, env, "orders", models.Order{Id: "o-cm1", Items: []models.OrderItem{{Id: "i1", Product: models.Product{EnableFixedCost: true}}}})
	seedServiceDoc(t, env, "orders", models.Order{Id: "o-cm2", Items: []models.OrderItem{{Id: "i2"}}})
	svc := newOrderSvc(env, "average")

	orders, _, err := svc.GetOrders(GetOrdersParameters{PageNumber: 1, PageSize: 10, FilterIsPaid: -1, IsPayLater: -1})
	require.NoError(t, err)
	require.Len(t, orders, 2)
}

func TestOrderService_SubmitOrder_PaymentsOnUnpaid(t *testing.T) {
	env := newServiceEnv(t)
	recipeSvc := newRecipeService(env)
	created, err := recipeSvc.InsertNew(models.Product{Name: "Pizza", Price: 10})
	require.NoError(t, err)

	svc := newOrderSvc(env, "average")
	_, err = svc.SubmitOrder(models.Order{
		IsPaid:   false,
		Payments: []models.OrderPayment{{Source: "Cash", Amount: 10}},
		Items:    []models.OrderItem{{Product: models.Product{Id: created.Id}, Quantity: 1}},
	})
	assert.Error(t, err)
}

func getOrderTips(t *testing.T, env *testutil.TestEnv, id string) float64 {
	t.Helper()
	var o models.Order
	require.NoError(t, env.Client.Database(env.Config.Databases[0].Database).Collection("orders").FindOne(context.Background(), bson.M{"id": id}).Decode(&o))
	return o.Tips
}
