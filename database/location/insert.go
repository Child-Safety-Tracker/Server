package location

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func DatabaseInsertLocation(database *pgx.Conn, ids []string) error {
	newLocationNums := 0 // Indicate the number of new locations
	// Fetch the locations from Apple server
	response, err := AppleServerLocations(ids, 7)

	if err != nil {
		log.Err(err).Msg("[Location] Failed to fetch the location from Apple Server")
		return err
	}

	// Fetch the latest location timestamp from the database
	var lastUpdateTimestamp int
	err = database.QueryRow(context.Background(), "SELECT max(\"DatePublished\") FROM \"DeviceLocation\";").Scan(&lastUpdateTimestamp)

	if err != nil {
		log.Err(err).Msg("[Database] Failed to fetch the latest location timestamp")
		return err
	}

	// The query string to insert the data into the database
	queryString := "INSERT INTO \"DeviceLocation\" (\"DeviceID\", \"DatePublished\", \"Description\", \"StatusCode\", \"Payload\") VALUES "

	// Add the value to be inserted
	for i := 0; i < len(response.Results); i++ {
		// Only location the location when there is new value available
		if response.Results[i].DatePublished > lastUpdateTimestamp {
			newLocationNums += 1
			queryString += "('" + response.Results[i].Id + "', " + fmt.Sprintf("%d", response.Results[i].DatePublished) + ", '" + response.Results[i].Description + "', " + fmt.Sprintf("%d", response.Results[i].StatusCode) + ", '" + response.Results[i].Payload + "')"

			// Add a comma to separate each insert data
			if i != len(response.Results)-1 {
				queryString += ","
			}
		}
	}

	// Trim the last trailing comma
	if newLocationNums != 0 {
		queryString = queryString[:len(queryString)-1]
	}

	if newLocationNums == 0 {
		// No location when there is no new location
		queryString = ""
		log.Warn().Msg("[Database] No new location to be updated")
	}

	// Execute the insert command
	_, err = database.Exec(context.Background(), queryString)

	if err != nil {
		log.Err(err).Msg("[Database] Error executing the location location query command")
		return err
	}

	log.Info().Msg("[Database] Successfully updated the location")
	return nil
}
