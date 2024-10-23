package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"os"
	"server/database/models"
	locationModels "server/location/models"
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

func GetUserInfo(database *pgx.Conn, userID string) (models.User, error) {
	user := models.User{}

	// Query the User from database and assign values into user variable
	err := database.QueryRow(context.Background(), "SELECT * FROM \"User\"").Scan(&user.UserID, &user.Username, &user.DeviceNums)
	if err != nil {
		log.Error().Err(err).Msg("[Database] Failed to get user info")
		return models.User{}, err
	}

	return user, nil
}

func GetDevicesInfo(database *pgx.Conn, userID string) ([]models.Device, error) {
	// One user can have many devices
	var devices []models.Device
	tempDevice := models.Device{}

	// Query the device(s) from database
	rows, err := database.Query(context.Background(), "SELECT * FROM \"Device\" WHERE \"UserID\" = $1", userID)
	if err != nil {
		log.Err(err).Msg("[Database] Failed to get devices info")
		return devices, err
	}

	// Assign the returned database query value to the array of Device
	for rows.Next() {
		err := rows.Scan(&tempDevice.DeviceID, &tempDevice.UserID, &tempDevice.PrivateKey)
		if err != nil {
			log.Err(err).Msg("[Database] Failed to scan device values.")
		}

		devices = append(devices, tempDevice)
	}

	return devices, nil
}

func InsertLocation(database *pgx.Conn, fetchedLocation locationModels.DecryptedLocationResult) error {
	// Assign the fetchedLocation to the Database model
	deviceLocation := models.DeviceLocation{
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
