package user

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	databaseModels "server/models/database"
)

func GetUserInfo(database *pgxpool.Pool, userID string) (databaseModels.User, error) {
	user := databaseModels.User{}

	// Query the User from database and assign values into user variable
	err := database.QueryRow(context.Background(), "SELECT * FROM \"User\" WHERE \"UserID\"=$1", userID).Scan(&user.UserID, &user.UserName, &user.DeviceNums)
	if err != nil {
		return databaseModels.User{}, err
	}

	return user, nil
}
