package database

import "backend/pkg/types"

// Database interface defines the methods that Database implementations will need to have
// Iterface is used although only one implementation is used so that we can mock it
type Database interface {
	GetDevices() ([]types.Device, error)
	InsertDevice(types.Device) error
	DeviceExistWithNameAndIP(string, string) (bool, error)
	DeviceFromName(string) (string, error)
}
