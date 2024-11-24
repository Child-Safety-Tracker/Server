package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"os"
	databaseModels "server/models/database"
	"server/models/location"
)

func DatabaseConnect() *pgx.Conn {
	database, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("[Database] Failed to connect")
	} else {
		log.Info().Msg("[Database] Connected to database")
	}

	return database
}

func InsertLocation(database *pgx.Conn, fetchedLocation location.DecryptedLocationResult) error {
	// Assign the fetchedLocation to the Database model
	deviceLocation := databaseModels.DeviceLocation{
		DeviceID:      fetchedLocation.Id,
		DatePublished: fetchedLocation.DatePublished,
		Description:   fetchedLocation.Description,
		StatusCode:    fetchedLocation.StatusCode,
		Latitude:      fetchedLocation.Payload.Latitude,
		Longitude:     fetchedLocation.Payload.Longitude,
		Confidence:    fetchedLocation.Payload.Confidence,
	}

	// Insert the device location into the database
	_, err := database.Exec(context.Background(), ""+
		"INSERT INTO \"DeviceLocation\""+"(\"DeviceID\", \"DatePublished\", \"Description\", \"StatusCode\", \"Latitude\", \"Longitude\", \"Confidence\") "+"VALUES ($1, $2, $3, $4, $5, $6, $7);", deviceLocation.DeviceID, deviceLocation.DatePublished, deviceLocation.Description, deviceLocation.StatusCode, deviceLocation.Latitude, deviceLocation.Longitude, deviceLocation.Confidence)
	if err != nil {
		log.Err(err).Msg("[Database] Failed to insert location")
		return err
	}

	// Print success message
	log.Info().Msg("[Database] Successfully inserted location")
	return nil
}
