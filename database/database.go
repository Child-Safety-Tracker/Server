package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func DatabaseConnect(connectionString string) *sql.DB {
	database, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal().Err(err).Msg("[Database] Failed to connect")
	} else {
		log.Info().Msg("[Database] Connected to database")
	}

	return database
}
