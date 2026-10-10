package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateProductImageHandler_BadMultipart(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "recipes", models.Product{Id: "p-1", Name: "Pizza"})

	router := mux.NewRouter()
	router.HandleFunc("/api/products/{id}/image", UpdateProductImage(env.Config, env.Logger)).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/products/p-1/image", "not multipart")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUpdateProductImageHandler_CreateFileError(t *testing.T) {
	env := newHandlerEnv(t)
	root := t.TempDir()
	// Point the uploads path at a regular file so os.Create under it fails.
	uploadsFile := filepath.Join(root, "uploads")
	require.NoError(t, os.WriteFile(uploadsFile, []byte("x"), 0o644))
	env.Config.UploadsPath = uploadsFile
	seed(t, env, "recipes", models.Product{Id: "p-1", Name: "Pizza"})

	router := mux.NewRouter()
	router.HandleFunc("/api/products/{id}/image", UpdateProductImage(env.Config, env.Logger)).Methods(http.MethodPost)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("image", "pic.png")
	require.NoError(t, err)
	_, err = part.Write([]byte("bytes"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req := newRequestWithBody(http.MethodPost, "/api/products/p-1/image", buf.Bytes())
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := serve(t, router, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetAvailableLanguagesHandler_MalformedJSON(t *testing.T) {
	dir := t.TempDir()
	langDir := filepath.Join(dir, "assets", "core", "languages")
	require.NoError(t, os.MkdirAll(filepath.Join(langDir, "a_dir"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(langDir, "b.json"), []byte("{not json"), 0o644))
	t.Chdir(dir)

	env := newHandlerEnv(t)
	router := mux.NewRouter()
	router.HandleFunc("/api/languages", GetAvailableLanguages(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/languages", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestProductHandlerErrorPaths(t *testing.T) {
	env := newHandlerEnv(t)

	router := mux.NewRouter()
	router.HandleFunc("/api/products/{id}", DeleteProduct(env.Config, env.Logger)).Methods(http.MethodDelete)
	router.HandleFunc("/api/products/{id}", UpdateProduct(env.Config, env.Logger)).Methods(http.MethodPatch)
	router.HandleFunc("/api/products/{id}/tree", GetRecipeTree(env.Config, env.Logger)).Methods(http.MethodGet)

	// Product image is set but the file no longer exists: delete fails on cleanup.
	seed(t, env, "recipes", models.Product{Id: "p-img", Name: "Pizza", ImageURL: "does-not-exist.png"})
	env.Config.UploadsPath = t.TempDir()
	rec := doRequest(t, router, http.MethodDelete, "/api/products/p-img", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// Update a missing product: update is a no-op, then get fails.
	rec = doRequest(t, router, http.MethodPatch, "/api/products/missing",
		jsonBody(t, map[string]interface{}{"data": models.Product{Name: "x"}}))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// Recipe tree for a missing product.
	rec = doRequest(t, router, http.MethodGet, "/api/products/missing/tree", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetSettingsHandler_Error(t *testing.T) {
	env := newHandlerEnv(t)
	require.NoError(t, env.Client.Database(env.Config.Databases[0].Database).Collection("settings").Drop(t.Context()))

	router := mux.NewRouter()
	router.HandleFunc("/api/settings", GetSettings(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/settings", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteEntryHandler_Error(t *testing.T) {
	env := newHandlerEnv(t)
	router := mux.NewRouter()
	router.HandleFunc("/api/materials/{material_id}/entries/{entry_id}", DeleteEntry(env.Config, env.Logger)).Methods(http.MethodDelete)

	rec := doRequest(t, router, http.MethodDelete, "/api/materials/missing/entries/e1", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetMaterialsHandler_DefaultPagination(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "materials", models.Material{Id: "m-1", Name: "Flour", Entries: []models.MaterialEntry{{Id: "e1", Quantity: 1}}})

	router := mux.NewRouter()
	router.HandleFunc("/api/materials", GetMaterials(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/materials", "")
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestFinishOrderHandler_MissingOrder(t *testing.T) {
	env := newHandlerEnv(t)
	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}/finish", FinishOrder(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/orders/missing/finish", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestSubmitOrderHandler_PaymentMismatch(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "recipes", models.Product{Id: "r-pay", Name: "Pizza", Price: 10})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/submit", SubmitOrder(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/orders/submit",
		jsonBody(t, map[string]interface{}{"data": models.Order{
			Items: []models.OrderItem{{Product: models.Product{Id: "r-pay"}, Quantity: 1}},
		}}))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetCategoriesHandler_DanglingProduct(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "categories", models.Category{Id: "c-1", Name: "Mains", Products: []models.Product{{Id: "missing"}}})

	router := mux.NewRouter()
	router.HandleFunc("/api/categories", GetCategories(env.Config, env.Logger)).Methods(http.MethodGet)

	// The service filters dangling product references, so the response is OK.
	rec := doRequest(t, router, http.MethodGet, "/api/categories", "")
	assert.Equal(t, http.StatusOK, rec.Code)
}
