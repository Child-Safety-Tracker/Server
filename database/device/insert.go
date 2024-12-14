package device

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"server/models/database"
)

func DatabaseInsertDevice(database *pgxpool.Pool, value database.Device) error {
	_, err := database.Exec(context.Background(), "INSERT INTO \"Device\"(\"DeviceID\", \"UserID\", \"PrivateKey\", \"Enabled\") VALUES ($1, $2, $3, $4);", value.DeviceID, value.UserID, value.PrivateKey, value.Enabled)

	if err != nil {
		log.Err(err).Msg("[Database] Failed inserting the device")
		return err
	}

	return nil
}
