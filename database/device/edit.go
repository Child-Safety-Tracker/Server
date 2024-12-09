package device

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func DatabaseSetDeviceStatus(database *pgx.Conn, deviceId string, status bool) error {

	// Convert to string for database query
	var convertedStatus string
	if status {
		convertedStatus = "true"
	} else {
		convertedStatus = "false"
	}

	_, err := database.Exec(context.Background(), "UPDATE \"Device\" SET \"Enabled\"=$1 WHERE \"DeviceID\"=$2", convertedStatus, deviceId)

	if err != nil {
		log.Err(err).Msg("[Database] Failed to set device status")
		return err
	}

	return nil
}
