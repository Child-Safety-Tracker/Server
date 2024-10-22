package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"os"
	"server/database/models"
)

func DatabaseConnect() *pgx.Conn {
	database, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("[Database] Failed to connect")
	} else {
		log.Info().Msg("[Database] Connected to database")
	}

	return database
}

func GetUserInfo(database *pgx.Conn, userID string) (models.User, error) {
	user := models.User{}

	// Query the User from database and assign values into user variable
	err := database.QueryRow(context.Background(), "SELECT * FROM \"User\"").Scan(&user.UserID, &user.Username, &user.DeviceNums)
	if err != nil {
		log.Error().Err(err).Msg("[Database] Failed to get user info")
		return models.User{}, err
	}

	return user, nil
}
