package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertNewProductHandler(t *testing.T) {
	env := newHandlerEnv(t)
	router := mux.NewRouter()
	router.HandleFunc("/api/products", InesrtNewProduct(env.Config, env.Logger)).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/products",
		jsonBody(t, map[string]interface{}{"data": models.Product{Name: "Pizza", Price: 12.5}}))
	require.Equal(t, http.StatusCreated, rec.Code)

	var p models.Product
	decodeData(t, rec, &p)
	assert.Equal(t, "Pizza", p.Name)
	assert.Equal(t, 12.5, p.Price)
	assert.NotEmpty(t, p.Id)
}

func TestInsertNewProductHandler_BadJSON(t *testing.T) {
	env := newHandlerEnv(t)
	router := mux.NewRouter()
	router.HandleFunc("/api/products", InesrtNewProduct(env.Config, env.Logger)).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/products", "{not json")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetProductHandler(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "recipes", models.Product{Id: "p-1", Name: "Pizza", Price: 12.5})

	router := mux.NewRouter()
	router.HandleFunc("/api/products/{id}", GetProduct(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/products/p-1", "")
	require.Equal(t, http.StatusOK, rec.Code)

	var p models.Product
	decodeData(t, rec, &p)
	assert.Equal(t, "Pizza", p.Name)
}

func TestGetProductHandler_NotFound(t *testing.T) {
	env := newHandlerEnv(t)
	router := mux.NewRouter()
	router.HandleFunc("/api/products/{id}", GetProduct(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/products/missing", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetProductsHandler(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "recipes",
		models.Product{Id: "p-1", Name: "Pizza", Price: 12.5},
		models.Product{Id: "p-2", Name: "Burger", Price: 8},
	)

	router := mux.NewRouter()
	router.HandleFunc("/api/products", GetProducts(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/products?page[number]=1&page[size]=10", "")
	require.Equal(t, http.StatusOK, rec.Code)

	var products []models.Product
	decodeData(t, rec, &products)
	assert.Len(t, products, 2)
	assert.Equal(t, 2, decodeMeta(t, rec).TotalRecords)

	// default pagination path
	rec = doRequest(t, router, http.MethodGet, "/api/products", "")
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateProductHandler(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "recipes", models.Product{Id: "p-1", Name: "Pizza", Price: 12.5})

	router := mux.NewRouter()
	router.HandleFunc("/api/products/{id}", UpdateProduct(env.Config, env.Logger)).Methods(http.MethodPatch)

	rec := doRequest(t, router, http.MethodPatch, "/api/products/p-1",
		jsonBody(t, map[string]interface{}{"data": models.Product{Name: "Pizza XL", Price: 15}}))
	require.Equal(t, http.StatusOK, rec.Code)

	var p models.Product
	decodeData(t, rec, &p)
	assert.Equal(t, "Pizza XL", p.Name)
	assert.Equal(t, 15.0, p.Price)
}

func TestUpdateProductHandler_BadJSON(t *testing.T) {
	env := newHandlerEnv(t)
	router := mux.NewRouter()
	router.HandleFunc("/api/products/{id}", UpdateProduct(env.Config, env.Logger)).Methods(http.MethodPatch)

	rec := doRequest(t, router, http.MethodPatch, "/api/products/p-1", "{")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteProductHandler(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "recipes", models.Product{Id: "p-1", Name: "Pizza"})

	router := mux.NewRouter()
	router.HandleFunc("/api/products/{id}", DeleteProduct(env.Config, env.Logger)).Methods(http.MethodDelete)

	rec := doRequest(t, router, http.MethodDelete, "/api/products/p-1", "")
	assert.Equal(t, http.StatusNoContent, rec.Code)

	count := countRecipes(t, env)
	assert.Equal(t, int64(0), count)
}

func TestDeleteProductHandler_NotFound(t *testing.T) {
	env := newHandlerEnv(t)
	router := mux.NewRouter()
	router.HandleFunc("/api/products/{id}", DeleteProduct(env.Config, env.Logger)).Methods(http.MethodDelete)

	rec := doRequest(t, router, http.MethodDelete, "/api/products/missing", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetRecipeTreeHandler(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "recipes", models.Product{Id: "r-1", Name: "Pizza", Price: 12, Quantity: 1})

	router := mux.NewRouter()
	router.HandleFunc("/api/products/{id}/tree", GetRecipeTree(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/products/r-1/tree", "")
	require.Equal(t, http.StatusOK, rec.Code)

	var tree models.Product
	decodeData(t, rec, &tree)
	assert.Equal(t, "r-1", tree.Id)
	assert.Equal(t, "Pizza", tree.Name)
}

func TestGetRecipeAvailabilityHandler(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "recipes", models.Product{Id: "r-1", Name: "Pizza", Ready: 5})

	router := mux.NewRouter()
	router.HandleFunc("/api/products/availability", GetRecipeAvailability(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/products/availability?ids=r-1", "")
	require.Equal(t, http.StatusOK, rec.Code)

	rec = doRequest(t, router, http.MethodGet, "/api/products/availability", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateProductImageHandler(t *testing.T) {
	env := newHandlerEnv(t)
	env.Config.UploadsPath = t.TempDir()
	seed(t, env, "recipes", models.Product{Id: "p-1", Name: "Pizza"})

	router := mux.NewRouter()
	router.HandleFunc("/api/products/{id}/image", UpdateProductImage(env.Config, env.Logger)).Methods(http.MethodPost)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("image", "pic.png")
	require.NoError(t, err)
	_, err = part.Write([]byte("fake-image-bytes"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req := newRequestWithBody(http.MethodPost, "/api/products/p-1/image", buf.Bytes())
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := serve(t, router, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var p models.Product
	require.NoError(t, env.Client.Database(env.Config.Databases[0].Database).
		Collection("recipes").FindOne(t.Context(), map[string]any{"id": "p-1"}).Decode(&p))
	assert.NotEmpty(t, p.ImageURL)
}
