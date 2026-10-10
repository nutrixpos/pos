package handlers

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomerHandlers(t *testing.T) {
	env := newHandlerEnv(t)
	settings := testHandlerSettings()

	router := mux.NewRouter()
	router.HandleFunc("/api/customers", AddCustomer(env.Config, env.Logger)).Methods(http.MethodPost)
	router.HandleFunc("/api/customers", GetCustomers(env.Config, env.Logger, settings)).Methods(http.MethodGet)
	router.HandleFunc("/api/customers/{id}", GetCustomer(env.Config, env.Logger)).Methods(http.MethodGet)
	router.HandleFunc("/api/customers/{id}", UpdateCustomer(env.Config, env.Logger)).Methods(http.MethodPatch)
	router.HandleFunc("/api/customers/{id}", DeleteCustomer(env.Config, env.Logger, settings)).Methods(http.MethodDelete)

	rec := doRequest(t, router, http.MethodPost, "/api/customers",
		jsonBody(t, map[string]interface{}{"data": models.Customer{Name: "Alice", Phone: "123", Address: "X"}}))
	require.Equal(t, http.StatusCreated, rec.Code)
	var created models.Customer
	decodeData(t, rec, &created)
	require.NotEmpty(t, created.Id)

	rec = doRequest(t, router, http.MethodGet, "/api/customers?page[number]=1&page[size]=10", "")
	require.Equal(t, http.StatusOK, rec.Code)
	var list []models.Customer
	decodeData(t, rec, &list)
	assert.Len(t, list, 1)

	rec = doRequest(t, router, http.MethodGet, "/api/customers/"+created.Id, "")
	require.Equal(t, http.StatusOK, rec.Code)
	var got models.Customer
	decodeData(t, rec, &got)
	assert.Equal(t, "Alice", got.Name)

	rec = doRequest(t, router, http.MethodPatch, "/api/customers/"+created.Id,
		jsonBody(t, map[string]interface{}{"data": models.Customer{Name: "Alice2", Phone: "456"}}))
	require.Equal(t, http.StatusOK, rec.Code)
	decodeData(t, rec, &got)
	assert.Equal(t, "Alice2", got.Name)

	rec = doRequest(t, router, http.MethodDelete, "/api/customers/"+created.Id, "")
	assert.Equal(t, http.StatusNoContent, rec.Code)

	rec = doRequest(t, router, http.MethodGet, "/api/customers/missing", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	rec = doRequest(t, router, http.MethodPost, "/api/customers", "{")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = doRequest(t, router, http.MethodPatch, "/api/customers/"+created.Id, "{")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCategoryHandlers(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "recipes", models.Product{Id: "p-1", Name: "Pizza"})

	router := mux.NewRouter()
	router.HandleFunc("/api/categories", InsertCategory(env.Config, env.Logger)).Methods(http.MethodPost)
	router.HandleFunc("/api/categories", GetCategories(env.Config, env.Logger)).Methods(http.MethodGet)
	router.HandleFunc("/api/categories", UpdateCategory(env.Config, env.Logger)).Methods(http.MethodPatch)
	router.HandleFunc("/api/categories/{id}", DeleteCategory(env.Config, env.Logger)).Methods(http.MethodDelete)

	rec := doRequest(t, router, http.MethodPost, "/api/categories",
		jsonBody(t, map[string]interface{}{"data": models.Category{Name: "Mains"}}))
	require.Equal(t, http.StatusCreated, rec.Code)

	rec = doRequest(t, router, http.MethodGet, "/api/categories?page[number]=1&page[size]=10", "")
	require.Equal(t, http.StatusOK, rec.Code)
	var list []models.Category
	decodeData(t, rec, &list)
	require.Len(t, list, 1)
	require.NotEmpty(t, list[0].Id)

	rec = doRequest(t, router, http.MethodPatch, "/api/categories",
		jsonBody(t, map[string]interface{}{"data": models.Category{Id: list[0].Id, Name: "Mains2", Products: []models.Product{{Id: "p-1"}}}}))
	require.Equal(t, http.StatusCreated, rec.Code)

	rec = doRequest(t, router, http.MethodDelete, "/api/categories/"+list[0].Id, "")
	require.Equal(t, http.StatusOK, rec.Code)

	rec = doRequest(t, router, http.MethodPost, "/api/categories", "{")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSettingsHandlers(t *testing.T) {
	env := newHandlerEnv(t)

	router := mux.NewRouter()
	router.HandleFunc("/api/settings", UpdateSettings(env.Config, env.Logger)).Methods(http.MethodPost)
	router.HandleFunc("/api/settings", GetSettings(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/settings", "")
	require.Equal(t, http.StatusOK, rec.Code)
	var got models.Settings
	decodeData(t, rec, &got)
	assert.Equal(t, "A", got.Orders.Queues[0].Prefix)

	rec = doRequest(t, router, http.MethodPost, "/api/settings",
		jsonBody(t, map[string]interface{}{"data": models.Settings{Language: models.LanguageSettings{Code: "ar", Language: "Arabic"}}}))
	require.Equal(t, http.StatusNoContent, rec.Code)

	rec = doRequest(t, router, http.MethodGet, "/api/settings", "")
	decodeData(t, rec, &got)
	assert.Equal(t, "ar", got.Language.Code)

	rec = doRequest(t, router, http.MethodPost, "/api/settings", "{")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSalesHandlers(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "sales", models.SalesPerDay{
		Date:       "2024-01-01",
		Costs:      5,
		TotalSales: 10,
		Orders: []models.SalesPerDayOrder{{
			Id:    "o-1",
			Order: models.Order{Id: "o-1", DisplayId: "A-1", Cost: 5, SalePrice: 10, Payments: []models.OrderPayment{{Source: "Cash", Amount: 10}}},
		}},
	})

	router := mux.NewRouter()
	router.HandleFunc("/api/sales", GetSalesPerDay(env.Config, env.Logger)).Methods(http.MethodGet)
	router.HandleFunc("/api/sales/export", ExportSalesCSV(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/sales?page[number]=1&page[size]=10", "")
	require.Equal(t, http.StatusOK, rec.Code)
	var days []models.SalesPerDay
	decodeData(t, rec, &days)
	require.Len(t, days, 1)
	assert.Equal(t, 1, decodeMeta(t, rec).TotalRecords)

	rec = doRequest(t, router, http.MethodGet, "/api/sales/export", "")
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "A-1")
	assert.Contains(t, rec.Body.String(), "Cash: 10.00")

	// Defaults path (no query params).
	rec = doRequest(t, router, http.MethodGet, "/api/sales", "")
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestLanguageHandlers(t *testing.T) {
	chdirRepoRoot(t)
	env := newHandlerEnv(t)

	router := mux.NewRouter()
	router.HandleFunc("/api/languages/{code}", GetLanguage(env.Config, env.Logger)).Methods(http.MethodGet)
	router.HandleFunc("/api/languages", GetAvailableLanguages(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/languages/en", "")
	require.Equal(t, http.StatusOK, rec.Code)

	rec = doRequest(t, router, http.MethodGet, "/api/languages", "")
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"en"`)
}
