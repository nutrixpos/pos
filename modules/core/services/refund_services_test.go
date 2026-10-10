package services

import (
	"testing"

	"github.com/nutrixpos/pos/internal/testutil"
	"github.com/nutrixpos/pos/modules/core/dto"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/require"
)

func seedRefundOrder(t *testing.T, env *testutil.TestEnv, orderID, itemID, productID string) {
	t.Helper()
	seedServiceDoc(t, env, "orders", models.Order{
		Id: orderID,
		Items: []models.OrderItem{{
			Id:        itemID,
			Product:   models.Product{Id: productID, Name: "Pizza"},
			Quantity:  2,
			SalePrice: 20,
			Materials: []models.OrderItemMaterial{{
				Material: models.Material{Id: "m-refund", Name: "Flour"},
				Entry:    models.MaterialEntry{Id: "e-refund"},
				Quantity: 1,
			}},
		}},
	})
}

func TestOrderService_RefundItem_MissingOrder(t *testing.T) {
	env := newServiceEnv(t)
	svc := newOrderSvc(env, "average")
	require.Error(t, svc.RefundItem(dto.OrderItemRefundRequest{OrderId: "missing", ItemId: "i"}, "u"))
}

func TestOrderService_RefundItem_Destinations(t *testing.T) {
	for _, dest := range []string{
		dto.DTOOrderItemRefundDestination_Inventory,
		dto.DTOOrderItemRefundDestination_Waste,
		dto.DTOOrderItemRefundDestination_Disposals,
	} {
		t.Run(dest, func(t *testing.T) {
			env := newServiceEnv(t)
			seedRefundOrder(t, env, "o-ref", "i-1", "p-1")
			svc := newOrderSvc(env, "average")

			err := svc.RefundItem(dto.OrderItemRefundRequest{
				OrderId:     "o-ref",
				ItemId:      "i-1",
				ProductId:   "p-1",
				Reason:      "customer",
				RefundValue: 5,
				Destination: dest,
			}, "u")
			require.NoError(t, err)
		})
	}
}

func TestOrderService_RefundItem_Custom(t *testing.T) {
	env := newServiceEnv(t)
	insertMaterial(t, env, models.Material{
		Id:      "m-refund",
		Name:    "Flour",
		Entries: []models.MaterialEntry{{Id: "e-refund", Quantity: 5, PurchaseQuantity: 5, PurchasePrice: 10}},
	})
	seedRefundOrder(t, env, "o-ref", "i-1", "p-1")
	svc := newOrderSvc(env, "average")

	err := svc.RefundItem(dto.OrderItemRefundRequest{
		OrderId:     "o-ref",
		ItemId:      "i-1",
		ProductId:   "p-1",
		Reason:      "customer",
		RefundValue: 5,
		Destination: dto.DTOOrderItemRefundDestination_Custom,
		MaterialRefunds: []dto.OrderItemRefundMaterialDTO{{
			MaterialId:         "m-refund",
			EntryId:            "e-refund",
			InventoryReturnQty: 2,
			DisposeQty:         1,
			WasteQty:           1,
		}},
		ProductAdd: []dto.OrderItemRefundProductAddDTO{{ProductId: "p-1", Quantity: 1}},
	}, "u")
	require.NoError(t, err)
}
