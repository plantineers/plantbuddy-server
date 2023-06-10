// Author: Yannick Kirschen
package controller

// ControllerRepository provides access to controller metadata.
type ControllerRepository interface {
	// GetAllUUIDs returns all UUIDs of all controllers.
	// Caution: This method does not use a transaction.
	GetAllUUIDs() ([]string, error)

	// GetByUUID returns the controller with the given UUID.
	// Caution: This method does not use a transaction.
	GetByUUID(uuid string) (*Controller, error)
}
