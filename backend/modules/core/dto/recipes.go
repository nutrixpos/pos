package dto

type RecipeAvailability struct {
	RecipeId              string             `json:"recipe_id"`
	Available             float64            `json:"available"`
	Ready                 float64            `json:"ready"`
	ComponentRequirements map[string]float64 `json:"component_requirements"`
}

type RecipeTree struct {
	RecipeId   string                `json:"recipe_id"`
	RecipeName string                `json:"recipe_name"`
	Components []PrepareItemResponse `json:"components"`
	SubRecipes []RecipeTree          `json:"sub_recipes"`
}
