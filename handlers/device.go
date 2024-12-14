package handlers

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"server/database/device"
	"server/models/database"
	"server/models/request"
)

func GetDevice(echoContext echo.Context, db *pgxpool.Pool) error {
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

func SetDeviceStatus(echoContext echo.Context, db *pgxpool.Pool) error {
	var requestBody request.DeviceStatusEditRequest
	err := echoContext.Bind(&requestBody)

	if err != nil {
		msg := "[Server] Failed to bind the request body"
		log.Err(err).Msg(msg)
		return echo.NewHTTPError(http.StatusInternalServerError, msg)
	}

	// Bad request
	if len(requestBody.DeviceId) == 0 {
		msg := "[Device] Invalid request body"
		log.Error().Msg(msg)
		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}

	err = device.DatabaseSetDeviceStatus(db, requestBody.DeviceId, requestBody.Enabled)

	if err != nil {
		msg := "[Database] Failed to set device status"
		log.Error().Msg(msg)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return echoContext.JSON(http.StatusOK, "[Device] Successfully set device status")
}

func InsertDevice(echoContext echo.Context, db *pgxpool.Pool) error {
	var requestBody database.Device
	err := echoContext.Bind(&requestBody)

	if err != nil {
		msg := "[Device] Failed to bind the request body"
		log.Err(err).Msg(msg)
		return echo.NewHTTPError(http.StatusInternalServerError, msg)
	}

	if len(requestBody.DeviceID) == 0 || len(requestBody.UserID) == 0 || len(requestBody.PrivateKey) == 0 {
		msg := "[Device] Invalid request body"
		log.Error().Msg(msg)
		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}

	err = device.DatabaseInsertDevice(db, requestBody)

	if err != nil {
		msg := "[Database] Failed to insert the device"
		log.Error().Msg(msg)
		return echo.NewHTTPError(http.StatusInternalServerError, msg)
	}

	return echoContext.JSON(http.StatusOK, "[Device] Successfully inserted the device")
}
