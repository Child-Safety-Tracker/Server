package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"os"
)

func DatabaseConnect() *pgxpool.Pool {
	database, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("[Database] Failed to connect")
	} else {
		log.Info().Msg("[Database] Connected to database")
	}

	return database
}
