package handlers

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"server/database/location"
	"server/models/request"
)

func GetLocations(echoContext echo.Context, database *pgxpool.Pool) error {

	// Bind the request body
	var requestBody request.LocationRequest
	err := echoContext.Bind(&requestBody)

	if err != nil {
		err = fmt.Errorf("[Server] Failed to bind the request body")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Bad request
	if len(requestBody.PrivateKeys) == 0 || len(requestBody.Ids) == 0 {
		msg := "[Location] Invalid request body"
		log.Error().Msg(msg)
		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}

	returnLocation, err := location.FetchLocation(database, requestBody.Ids, requestBody.PrivateKeys)

	if err != nil {
		msg := "[Location] Failed to fetch and decrypt location"
		log.Error().Msg(msg)
		return echo.NewHTTPError(http.StatusInternalServerError, msg)
	}

	return echoContext.JSON(http.StatusOK, returnLocation)
}
