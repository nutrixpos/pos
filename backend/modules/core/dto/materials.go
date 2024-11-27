// Package dto contains structs for data transfer objects (DTOs) used
// in the core module of nutrix.
package dto

import (
	"github.com/elmawardy/nutrix/backend/modules/core/models"
)

// MaterialEditRequest is a DTO used for editing material details.
type MaterialEditRequest struct {
	Material models.Material
}
