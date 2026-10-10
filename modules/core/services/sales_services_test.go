package services

import (
	"testing"

	"github.com/nutrixpos/pos/modules/core/dto"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSalesService_AddOrderItemToDayRefund(t *testing.T) {
	env := newServiceEnv(t)
	seedServiceDoc(t, env, "orders", models.Order{Id: "o-r", Items: []models.OrderItem{{Id: "i-1", Cost: 5}}})
	svc := SalesService{Logger: env.Logger, Config: env.Config}

	req := dto.OrderItemRefundRequest{OrderId: "o-r", ItemId: "i-1", Reason: "x", RefundValue: 5, Destination: "waste"}
	require.NoError(t, svc.AddOrderItemToDayRefund(req, "u")) // creates the day
	require.NoError(t, svc.AddOrderItemToDayRefund(req, "u")) // pushes into the existing day

	err := svc.AddOrderItemToDayRefund(dto.OrderItemRefundRequest{OrderId: "missing", ItemId: "i"}, "u")
	assert.Error(t, err)
}

func TestPurchaseOrderService_Reads(t *testing.T) {
	env := newServiceEnv(t)
	insertMaterialForPO(t, env, "mat-read")

	svc := newPOService(env)
	po, err := svc.CreatePurchaseOrder(makeAutoPO("mat-read", 5), "user-1")
	require.NoError(t, err)

	list, total, err := svc.GetPurchaseOrders(GetPurchaseOrdersParams{PageNumber: 1, PageSize: 10, Search: "test", Status: models.PurchaseOrderStatusReceived})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, list, 1)

	got, err := svc.GetPurchaseOrder(po.Id)
	require.NoError(t, err)
	assert.Equal(t, po.Id, got.Id)

	_, err = svc.GetPurchaseOrder("missing")
	assert.Error(t, err)

	grns, grnTotal, err := svc.GetGRNs(GetGRNsParams{PageNumber: 1, PageSize: 10, PurchaseOrderId: po.Id})
	require.NoError(t, err)
	assert.Equal(t, int64(1), grnTotal)
	require.Len(t, grns, 1)

	grn, err := svc.GetGRN(grns[0].Id)
	require.NoError(t, err)
	assert.Equal(t, po.Id, grn.PurchaseOrderId)

	_, err = svc.GetGRN("missing")
	assert.Error(t, err)
}

func TestSalesService(t *testing.T) {
	env := newServiceEnv(t)
	svc := SalesService{Logger: env.Logger, Config: env.Config}

	order := models.Order{Id: "o-sale", DisplayId: "A-1", Cost: 5, SalePrice: 10}

	// Creates today's sales document.
	require.NoError(t, svc.AddOrderToSalesDay(order, []models.ItemCost{}, "u"))
	// Second call exercises the push/update branch.
	require.NoError(t, svc.AddOrderToSalesDay(models.Order{Id: "o-sale-2", DisplayId: "A-2", Cost: 2, SalePrice: 4}, []models.ItemCost{}, "u"))

	days, total, err := svc.GetSalesPerday(1, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	require.Len(t, days, 1)
	assert.Len(t, days[0].Orders, 2)

	// Update the tips of an existing order.
	order.Tips = 3
	require.NoError(t, svc.SetOrderToSalesDay(order))

	// Updating an order that is not in any sales day is a no-op.
	require.NoError(t, svc.SetOrderToSalesDay(models.Order{Id: "not-there"}))
}
