package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/services"
)

func GetRecipeTree(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "GET, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "recipe id query string is required", http.StatusBadRequest)
			return
		}

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		tree, err := recipeService.GetRecipeTree(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tree); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetRecipeAvailability(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "GET, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		id := r.URL.Query().Get("ids")
		if id == "" {
			http.Error(w, "recipe ids comma separated is required as query string", http.StatusBadRequest)
			return
		}

		ids := strings.Split(id, `,`)

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		availabilities, err := recipeService.CheckRecipesAvailability(ids)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(availabilities); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}
