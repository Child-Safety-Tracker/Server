package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"server/database"
	"server/handlers"
	"time"
)

func main() {

	// Load environment variables from .env file
	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get project directory")
	}

	err = godotenv.Load(rootPath + "/.env")
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

	// -- Routes --
	echoInstance.GET("/", hello)

	// Location
	echoInstance.POST("/location", func(echoContext echo.Context) error {
		return handlers.GetLocations(echoContext, db)

	})

	// User
	// Using anonymous function to pass more arguments into handler function
	echoInstance.GET("/user", func(echoContext echo.Context) error {
		return handlers.GetUser(echoContext, db)
	})
	echoInstance.POST("/user/register", handlers.UserRegister)
	echoInstance.POST("/user/login", handlers.UserLogin)

	// Device
	echoInstance.GET("/device", func(echoContext echo.Context) error {
		return handlers.GetDevice(echoContext, db)
	})
	echoInstance.POST("/device/status", func(echoContext echo.Context) error {
		return handlers.SetDeviceStatus(echoContext, db)
	})
	echoInstance.POST("/device", func(echoContext echo.Context) error {
		return handlers.InsertDevice(echoContext, db)
	})

	// Start the server and logging result
	logger.Fatal().Err(echoInstance.Start(":1234")).Msg("[Server] Failed to start the server.")

	// Deference
	defer db.Close()
}

func hello(echoContext echo.Context) error {
	return echoContext.String(http.StatusOK, "Hello world")
}
