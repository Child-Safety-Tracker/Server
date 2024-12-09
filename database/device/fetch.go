package device

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	databaseModels "server/models/database"
)

func GetDevicesInfo(database *pgx.Conn, userID string) ([]databaseModels.Device, error) {
	// One user can have many devices
	var devices []databaseModels.Device
	var queriedRows pgx.Rows
	var err error
	tempDevice := databaseModels.Device{}

	// Query the device(s) from database
	// Query only for a user
	if len(userID) > 0 {
		queriedRows, err = database.Query(context.Background(), "SELECT * FROM \"Device\" WHERE \"UserID\" = $1", userID)
		if err != nil {
			log.Err(err).Msg("[Database] Failed to get devices info")
			return devices, err
		}
		// Query for all devices
	} else {
		queriedRows, err = database.Query(context.Background(), "SELECT * FROM \"Device\"")

		if err != nil {
			log.Err(err).Msg("[Database] Failed to get all of the devices")
			return devices, err
		}
	}

	// Assign the returned database query value to the array of Device
	for queriedRows.Next() {
		err := queriedRows.Scan(&tempDevice.DeviceID, &tempDevice.UserID, &tempDevice.PrivateKey, &tempDevice.Enabled)
		if err != nil {
			log.Err(err).Msg("[Database] Failed to scan device values.")
			return devices, err
		}

		devices = append(devices, tempDevice)
	}
	// Assign the result to an array of Device(s)
	for queriedRows.Next() {
		err := queriedRows.Scan(&tempDevice.DeviceID, &tempDevice.UserID, &tempDevice.PrivateKey)
		if err != nil {
			log.Err(err).Msg("[Database] Failed to scan device values.")
			return devices, err
		}
		devices = append(devices, tempDevice)
	}

	return devices, nil
}
