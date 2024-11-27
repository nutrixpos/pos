// Package dto contains structs for data transfer objects (DTOs) used
// in the core module of nutrix.
package dto

// ComponentQuantity is a DTO containing the ID of a component and the
// quantity of this component.
type ComponentQuantity struct {
	ComponentId string  `json:"component_id"` // The ID of the component.
	Quantity    float64 `json:"quantity"`     // The quantity of the component.
}
