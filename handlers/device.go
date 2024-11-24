package handlers

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"server/database/device"
)

func GetDevice(echoContext echo.Context, db *pgx.Conn) error {
	// Query parameter
	userID := echoContext.QueryParam("userId")

	// Query the device information
	result, err := device.GetDevicesInfo(db, userID)

	if err != nil {
		msg := "[Device] Failed query device information from the database"
		log.Error().Msg(msg)
		err = fmt.Errorf(msg)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return echoContext.JSON(http.StatusOK, result)
}
