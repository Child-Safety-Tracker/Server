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
	var fetchedLocationResult models.PostResponseBody
	var decryptedLocationResultValue models.DecryptedLocationResult
	var requestBody models.PostRequestBody

	// Bind the request body
	err := echoContext.Bind(&requestBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if len(requestBody.PrivateKey) == 0 || len(requestBody.Ids) == 0 || (requestBody.Days == 0) {
		print(requestBody.Days)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Fetch locations from Apple Server
	fetchedLocationResult, err = location.FetchLocation(URL, requestBody.Ids, requestBody.Days)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(fetchedLocationResult.Results) == 0 {
		return echoContext.JSON(http.StatusOK, models.PostResponseBody{Results: []models.LocationResult{}, StatusCode: "200"})
	} else {
		// Decrypt the newest fetched location
		decryptedLocationResultValue, err = decrypt.DecryptLocation(fetchedLocationResult.Results[0], requestBody.PrivateKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		// Return response
		return echoContext.JSON(http.StatusOK, decryptedLocationResultValue)
	}
}
