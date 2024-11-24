package handlers

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"server/database/user"
)

func GetUser(echoContext echo.Context, db *pgx.Conn) error {
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
