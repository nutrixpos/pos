package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/services"
)

func GetProductReadyNumber(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id query string is required", http.StatusBadRequest)
			return
		}

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		ready, err := recipeService.GetReadyNumber(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ready); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetRecipeTree(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

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
