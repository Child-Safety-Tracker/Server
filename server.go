package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"
)

func main() {
	// Init Echo instance
	echoInstance := echo.New()

	// Logger init with pretty format and timestamp enabled
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

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

	// Routes
	echoInstance.GET("/", hello)

	// Start the server and logging result
	logger.Fatal().AnErr("Error:", echoInstance.Start(":1234"))

}

func hello(echoContext echo.Context) error {
	return echoContext.String(http.StatusOK, "Hello world")
}
