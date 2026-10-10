package services

import (
	"net/http"
	"testing"
	"time"

	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeNotification struct {
	topics []string
}

func (f *fakeNotification) HandleHttpRequest(http.ResponseWriter, *http.Request) error { return nil }
func (f *fakeNotification) SendToTopic(topic, _ string) error {
	f.topics = append(f.topics, topic)
	return nil
}

func TestCustomersService_CRUD(t *testing.T) {
	env := newServiceEnv(t)
	svc := CustomersService{Logger: env.Logger, Config: env.Config}

	created, err := svc.InsertNew(models.Customer{Name: "Alice", Phone: "1"})
	require.NoError(t, err)
	require.NotEmpty(t, created.Id)

	got, err := svc.GetCustomer(created.Id)
	require.NoError(t, err)
	assert.Equal(t, "Alice", got.Name)

	list, count, err := svc.GetCustomers(GetCustomersParams{PageNumber: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, 1, count)

	updated, err := svc.UpdateCustomer(models.Customer{Name: "Alice2", Address: "X"}, created.Id)
	require.NoError(t, err)
	assert.Equal(t, "Alice2", updated.Name)

	require.NoError(t, svc.DeleteCustomer(created.Id))
	_, err = svc.GetCustomer(created.Id)
	assert.Error(t, err)
}

func TestCategoryService_CRUD(t *testing.T) {
	env := newServiceEnv(t)
	svc := CategoryService{Logger: env.Logger, Config: env.Config}
	_, err := (&RecipeService{Logger: env.Logger, Config: env.Config}).InsertNew(models.Product{Name: "Pizza"})
	require.NoError(t, err)

	require.NoError(t, svc.InsertCategory(models.Category{Name: "Mains"}))

	categories, err := svc.GetCategories(1, 10)
	require.NoError(t, err)
	require.Len(t, categories, 1)
	require.NotEmpty(t, categories[0].Id)

	_, err = svc.UpdateCategory(models.Category{Id: categories[0].Id, Name: "Mains2"})
	require.NoError(t, err)

	require.NoError(t, svc.DeleteCategory(categories[0].Id))
	categories, err = svc.GetCategories(1, 10)
	require.NoError(t, err)
	assert.Empty(t, categories)
}

func TestSettingsService_GetUpdate(t *testing.T) {
	env := newServiceEnv(t)
	svc := SettingsService{Logger: env.Logger, Config: env.Config}

	require.NoError(t, svc.UpdateSettings(models.Settings{Language: models.LanguageSettings{Code: "ar"}}))
	settings, err := svc.GetSettings()
	require.NoError(t, err)
	assert.Equal(t, "ar", settings.Language.Code)
}

func TestLogService_Queries(t *testing.T) {
	env := newServiceEnv(t)
	seedServiceDoc(t, env, "logs", map[string]interface{}{
		"id": "log-1", "type": models.LogTypeMaterialGRNReceive, "material_id": "m-1", "date": nowTime(),
	})
	seedServiceDoc(t, env, "logs", map[string]interface{}{
		"id": "log-2", "type": models.LogTypeOrderFinish, "order_id": "o-1",
	})
	seedServiceDoc(t, env, "logs", map[string]interface{}{
		"id": "log-3", "type": models.LogTypeOrderItemRefunded, "order_id": "o-1",
	})

	svc := LogService{Logger: env.Logger, Config: env.Config}

	logs, total, err := svc.GetMaterialLogs("m-1", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, logs, 1)

	salesLogs := svc.GetSalesLogs()
	assert.Len(t, salesLogs, 1)

	_, err = svc.GetOrderItemsRefundLogs([][]string{{"o-1"}})
	assert.NoError(t, err)
}

func TestLanguageService_GetLanguage(t *testing.T) {
	t.Chdir("../../../")
	env := newServiceEnv(t)
	svc := LanguageService{Logger: env.Logger, Config: env.Config}

	lang, err := svc.GetLanguage("en")
	require.NoError(t, err)
	assert.Equal(t, "en", lang.Code)

	// Unknown code returns the zero value without error.
	lang, err = svc.GetLanguage("zz")
	require.NoError(t, err)
	assert.Empty(t, lang.Code)
}

func TestDisposalService_CRUD(t *testing.T) {
	env := newServiceEnv(t)
	svc := DisposalService{Logger: env.Logger, Config: env.Config}

	materialDisposal := models.MaterialDisposal{
		Disposal:   models.Disposal{OrderId: "o-1", Type: models.TypeDisposalMaterial, Quantity: 2},
		MaterialId: "m-1",
		EntryId:    "e-1",
	}
	require.NoError(t, svc.AddMaterialDisposal(materialDisposal, "user"))

	productDisposal := models.ProductDisposal{
		Disposal: models.Disposal{OrderId: "o-1", Type: models.TypeDisposalProduct, Quantity: 1},
	}
	require.NoError(t, svc.AddProductDisposal(productDisposal, "user"))

	list, total, err := svc.GetDisposals(GetDisposalsParameters{PageNumber: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)

	_, err = svc.UpdateDisposal("x", materialDisposal)
	require.NoError(t, err)

	_, err = svc.UpdateDisposal("x", productDisposal)
	require.NoError(t, err)

	require.NoError(t, svc.DeleteDisposal("x"))
}

func TestCheckExpirationDates(t *testing.T) {
	env := newServiceEnv(t)
	env.Config.Databases[0].Name = env.Config.Databases[0].Database
	seedServiceDoc(t, env, "materials", models.Material{
		Id:   "m-exp",
		Name: "Milk",
		Entries: []models.MaterialEntry{{
			Id:             "e1",
			Quantity:       1,
			ExpirationDate: nowTime().Add(24 * time.Hour),
		}},
	})

	notifier := &fakeNotification{}
	CheckExpirationDates(env.Logger, env.Config, notifier)
	assert.Contains(t, notifier.topics, "expire_soon")
}
