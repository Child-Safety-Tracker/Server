package handlers

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"server/database/user"
)

func GetUser(echoContext echo.Context, db *pgxpool.Pool) error {
	// Query parameter
	userID := echoContext.QueryParam("userId")

	// Query the user information and put them into the user object
	result, err := user.GetUserInfo(db, userID)

	// Error response
	if err != nil {
		msg := "[User] Failed query user from the database"
		log.Error().Msg(msg)
		err = fmt.Errorf(msg)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return echoContext.JSON(http.StatusOK, result)
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserRegister(echoContext echo.Context) error {
	var requestBody UserCredentials
	err := echoContext.Bind(&requestBody)

	if err != nil {
		err = fmt.Errorf("[Register] Failed to bind the request body")
		log.Err(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(requestBody.Username) == 0 || len(requestBody.Password) == 0 {
		err = fmt.Errorf("[Register] Invalid user information")
		log.Err(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if requestBody.Username == "test" {
		err = fmt.Errorf("[Register] User already exist")
		log.Err(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return echoContext.JSON(http.StatusOK, "[Register] User registered successfully")
}
