package database

import "backend/pkg/types"

// Database interface defines the methods that Database implementations will need to have
// Iterface is used although only one implementation is used so that we can mock it
type Database interface {

	/*
		Devices management
	*/

	GetDevices() ([]types.Device, error)
	GetDeviceByUUID(string) (types.Device, error)
	InsertDevice(types.Device) error
	DeviceExistWithNameAndIP(string, string) (bool, error)
	DeviceIPFromName(string) (string, error)
	DeviceIPAndUUIDFromName(string) (string, string, error)
	DeleteDeviceFromUUID(string) error
	UpdateDevice(types.Device) error

	/*
		Messages and results management
	*/

	InsertMessage(types.MessageDB) error
	InsertResult(types.ResultDB) error
}
