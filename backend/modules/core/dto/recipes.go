package dto

type RecipeAvailability struct {
	RecipeId              string             `json:"recipe_id"`
	Available             float64            `json:"available"`
	Ready                 float64            `json:"ready"`
	ComponentRequirements map[string]float64 `json:"component_requirements"`
}
