package main

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"server/database"
	"server/database/device"
	"server/database/location"
	"server/models/response"
	"time"
)

func main() {
	var deviceIds []string
	var fetchedLocations response.LocationResponse

	// Load environment variables from .env file
	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get project directory")
	}

	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		log.Fatal().Err(err).Msg("[Server ]Error loading .env file")
	}

	// Logger init with pretty format and timestamp enabled
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
	// Set the logger to use globally
	log.Logger = logger

	// Connect to the database
	db := database.DatabaseConnect()

	// Get all available devices
	devices, err := device.GetDevicesInfo(db, "")
	if err != nil {
		log.Error().Msg("[Device] Error getting device info")
	}

	// Copy the device ids into another array
	for _, device := range devices {
		deviceIds = append(deviceIds, device.DeviceID)
	}

	// Fetch location from Apple server
	fetchedLocations, err = location.AppleServerLocations(deviceIds, 7)
	if err != nil {
		log.Error().Msg("[Location] Error fetching location from Apple server")
	}

	// Insert the new locations into the database
	err = location.DatabaseInsertLocation(db, fetchedLocations)
	if err != nil {
		log.Error().Msg("[Database] Error insert all locations")
	}

}
