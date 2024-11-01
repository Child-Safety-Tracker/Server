package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"server/database"
	"server/handlers"
	"server/location"
	"server/location/decrypt"
	"time"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("[Server ]Error loading .env file")
	}

	// Init Echo instance
	echoInstance := echo.New()

	// Logger init with pretty format and timestamp enabled
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
	// Set the logger to use globally
	log.Logger = logger

	// Middleware
	echoInstance.Use(middleware.Recover())
	// Logger Middleware with enabled log sections
	echoInstance.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogRemoteIP: true,
		LogValuesFunc: func(echoContext echo.Context, loggerValues middleware.RequestLoggerValues) error {
			logger.Info().Str("Address", loggerValues.RemoteIP).Str("URI", loggerValues.URI).Int("Status", loggerValues.Status).Msg("Request")
			return nil
		},
	}))

	// Connect to the database
	db := database.DatabaseConnect()

	// Routes
	echoInstance.GET("/", hello)
	echoInstance.POST("/location", handlers.GetLocations)
	fetchedLocation, err := location.FetchLocation(os.Getenv("APPLE_SERVER_WRAPPER_URL"), []string{"afirx1LlNk5vh7BnbGukU+L8o9E3pHhd/uogNOdmdv8="}, 5)
	decrypted, _ := decrypt.DecryptLocation(fetchedLocation.Results[0], "hUotVQIdoniIfacuUNHahmnNK98GRV6+kn+sOQ==")
	_ = database.InsertLocation(db, decrypted)

	// Start the server and logging result
	logger.Fatal().Err(echoInstance.Start(":1234")).Msg("[Server] Failed to start the server.")

	// Deference
	defer db.Close(context.Background())
}

func hello(echoContext echo.Context) error {
	return echoContext.String(http.StatusOK, "Hello world")
}
