package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"server/location"
	"server/location/decrypt"
	"server/location/models"
)

var URL string = "http://104.214.184.97:6176"

func GetLocations(echoContext echo.Context) error {
	fetchedLocationResult := models.PostResponseBody{}
	decryptedLocationResultValue := models.DecryptedLocationResult{}
	requestBody := models.PostRequestBody{}

	// Bind the request body
	err := echoContext.Bind(&requestBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Fetch locations from Apple Server
	fetchedLocationResult, err = location.FetchLocation(URL, requestBody.Ids, requestBody.Days)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Decrypt the newest fetched location
	if len(fetchedLocationResult.Results) == 0 {
		return echoContext.JSON(http.StatusOK, models.PostResponseBody{Results: []models.LocationResult{}, StatusCode: "200"})
	} else {
		decryptedLocationResultValue, err = decrypt.DecryptLocation(fetchedLocationResult.Results[0])
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		// Return response
		return echoContext.JSON(http.StatusOK, decryptedLocationResultValue)
	}
}
