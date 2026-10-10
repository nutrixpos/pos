package services

import (
	"testing"

	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisposalService_GetDisposals_Filters(t *testing.T) {
	env := newServiceEnv(t)
	seedServiceDoc(t, env, "disposals", models.Disposal{Id: "d-abc", Type: models.TypeDisposalMaterial, Quantity: 1})
	svc := DisposalService{Logger: env.Logger, Config: env.Config}

	list, _, err := svc.GetDisposals(GetDisposalsParameters{
		PageNumber:         1,
		PageSize:           10,
		DisposalIdContains: "d-abc",
		FilterState:        []string{"active", "!cancelled"},
	})
	require.NoError(t, err)
	_ = list
}

func TestMaterialService_GetMaterialEntries_Errors(t *testing.T) {
	env := newServiceEnv(t)
	insertMaterial(t, env, models.Material{Id: "mat-e", Entries: []models.MaterialEntry{{Id: "e1", Quantity: 1}}})
	svc := newMaterialService(env, models.Settings{})

	_, _, err := svc.GetMaterialEntries("missing", GetMaterialEntriesParams{PageNumber: 1, PageSize: 10})
	assert.Error(t, err)

	_, _, err = svc.GetMaterialEntries("mat-e", GetMaterialEntriesParams{PageNumber: 1, PageSize: -1})
	assert.Error(t, err)
}

func TestSeeder_SeedCategories_FreshNoProducts(t *testing.T) {
	env := newServiceEnv(t)
	svc := Seeder{Logger: env.Logger, Config: env.Config, Prompter: fakePrompter{confirm: false}}

	// No categories and no products: the prompter declines and the empty
	// category is still inserted.
	require.NoError(t, svc.SeedCategories())
}

func TestPurchaseOrderService_GetGRNs_Search(t *testing.T) {
	env := newServiceEnv(t)
	insertMaterialForPO(t, env, "mat-grn-s")
	svc := newPOService(env)
	_, err := svc.CreatePurchaseOrder(makeAutoPO("mat-grn-s", 2), "u")
	require.NoError(t, err)

	grns, _, err := svc.GetGRNs(GetGRNsParams{PageNumber: 1, PageSize: 10, Search: "PO"})
	require.NoError(t, err)
	assert.NotEmpty(t, grns)
}

func TestReceivePurchaseOrder_UnknownRequestedItem(t *testing.T) {
	env := newServiceEnv(t)
	insertMaterialForPO(t, env, "mat-unknown")
	svc := newPOService(env)
	po, err := svc.CreatePurchaseOrder(models.PurchaseOrder{
		Supplier: "s",
		Items:    []models.PurchaseOrderItem{{MaterialId: "mat-unknown", Quantity: 5, PurchasePrice: 1}},
	}, "u")
	require.NoError(t, err)

	// Requesting an item id that is not part of the order leaves nothing to receive.
	_, err = svc.ReceivePurchaseOrder(po.Id, "u", []ReceiveItem{{ItemId: "not-an-item", Quantity: 1}})
	assert.Error(t, err)
}
