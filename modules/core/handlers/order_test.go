package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/internal/testutil"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func insertOrderDoc(t *testing.T, env *testutil.TestEnv, order models.Order) {
	t.Helper()
	_, err := env.Client.Database(env.Config.Databases[0].Database).Collection("orders").
		InsertOne(context.Background(), order)
	require.NoError(t, err)
}

func findOrderDoc(t *testing.T, env *testutil.TestEnv, id string) models.Order {
	t.Helper()
	var o models.Order
	err := env.Client.Database(env.Config.Databases[0].Database).Collection("orders").
		FindOne(context.Background(), bson.M{"id": id}).Decode(&o)
	require.NoError(t, err)
	return o
}

func TestOrderAddTipHandler(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertOrderDoc(t, env, models.Order{Id: "order-1", Tips: 0})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{order_id}/addtips", OrderAddTip(env.Config, env.Logger, models.Settings{})).Methods("PATCH")

	rec := doRequest(t, router, http.MethodPatch, "/api/orders/order-1/addtips?tip_amount=5", "")
	require.Equal(t, http.StatusOK, rec.Code)

	stored := findOrderDoc(t, env, "order-1")
	assert.Equal(t, 5.0, stored.Tips)
}

func TestOrderAddTipHandler_MissingAmount(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{order_id}/addtips", OrderAddTip(env.Config, env.Logger, models.Settings{})).Methods("PATCH")

	rec := doRequest(t, router, http.MethodPatch, "/api/orders/order-1/addtips", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestOrderRemoveTipHandler(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertOrderDoc(t, env, models.Order{Id: "order-2", Tips: 10})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{order_id}/removetips", OrderRemoveTip(env.Config, env.Logger, models.Settings{})).Methods("PATCH")

	rec := doRequest(t, router, http.MethodPatch, "/api/orders/order-2/removetips?tip_amount=3", "")
	require.Equal(t, http.StatusOK, rec.Code)

	stored := findOrderDoc(t, env, "order-2")
	assert.Equal(t, 7.0, stored.Tips)
}
