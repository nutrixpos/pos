// Package dto contains Data Transfer Objects (DTOs) which are used to transfer data between application components.
// It is usually used for client-server communication.
package dto

// HttpComponent is a DTO containing the most important information about Material.
// It is used to return data from the API to the client.
type HttpComponent struct {
	Name     string  `json:"name"`
	Unit     string  `json:"unit"`
	Quantity float32 `json:"quantity"`
	Company  string  `json:"company"`
}

// GetComponentConsumeLogsRequest is a DTO used in the request body in the GET /material/consume_logs endpoint.
// It contains the material name.
type GetComponentConsumeLogsRequest struct {
	Name string `json:"name"`
}
