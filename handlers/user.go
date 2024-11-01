package handlers

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"server/database"
)

func GetUser(echoContext echo.Context, db *pgx.Conn) error {
	// Query parameter
	userID := echoContext.QueryParam("userId")

	// Query the user information and put them into the user object
	result, err := database.GetUserInfo(db, userID)

	// Error response
	if err != nil {
		log.Err(err).Msg("[User]")
		err = fmt.Errorf("[User] Failed query user from the database")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return echoContext.JSON(http.StatusOK, result)
}
