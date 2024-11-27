// RecipeAvailability is a DTO containing the id of a recipe and the availability
// of that recipe. The availability is a sum of the available and ready quantity.
// The component requirements are also included in this DTO.
package dto

// RecipeAvailability is a DTO containing the id of a recipe and the availability
// of that recipe. The availability is a sum of the available and ready quantity.
// The component requirements are also included in this DTO.
//
// The component requirements are a map of component id to the required quantity.
type RecipeAvailability struct {
	RecipeId              string             `json:"recipe_id"`
	Available             float64            `json:"available"`
	Ready                 float64            `json:"ready"`
	ComponentRequirements map[string]float64 `json:"component_requirements"`
}
