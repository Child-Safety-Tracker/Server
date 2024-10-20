package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Init Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start the server and logging result
	e.Logger.Fatal(e.Start(":1234"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world")
}
