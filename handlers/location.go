package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"server/location"
	"server/location/decrypt"
	locationModels "server/models/location"
	"server/models/request"
	"server/models/response"
)

func GetLocations(echoContext echo.Context) error {
	// Get URL from environment variables
	var URL string = os.Getenv("APPLE_SERVER_WRAPPER_URL")

	var fetchedLocationResult response.LocationResponse
	var decryptedLocationResultValue locationModels.DecryptedLocationResult
	var requestBody request.LocationRequest

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
		return echoContext.JSON(http.StatusOK, response.LocationResponse{Results: []locationModels.LocationResult{}, StatusCode: "200"})
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
