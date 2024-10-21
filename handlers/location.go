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

	err := echoContext.Bind(&requestBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	fetchedLocationResult, err = location.FetchLocation(URL, requestBody.Ids, requestBody.Days)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	decryptedLocationResultValue, err = decrypt.DecryptLocation(fetchedLocationResult.Results[0])
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return echoContext.JSON(http.StatusOK, decryptedLocationResultValue)
}
