package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPrinters(t *testing.T) {
	env := newHandlerEnv(t)

	router := mux.NewRouter()
	router.HandleFunc("/api/printers", GetPrinters(env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/printers", "")
	require.Equal(t, http.StatusOK, rec.Code)

	var response struct {
		Data struct {
			OS       string   `json:"os"`
			Printers []string `json:"printers"`
			Browsers []string `json:"browsers"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	assert.NotEmpty(t, response.Data.OS)
	assert.NotNil(t, response.Data.Printers)
	assert.NotNil(t, response.Data.Browsers)
}

func TestGetPrinters_ResponseContentType(t *testing.T) {
	env := newHandlerEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/printers", nil)
	rec := httptest.NewRecorder()
	GetPrinters(env.Logger)(rec, req)

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}
