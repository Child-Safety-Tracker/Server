package location

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"server/models/response"
)

func DatabaseInsertLocation(database *pgxpool.Pool, fetchedLocations response.LocationResponse) error {
	newLocationNums := 0 // Indicate the number of new locations

	// Fetch the latest location timestamp from the database
	var lastUpdateTimestamp int
	err := database.QueryRow(context.Background(), "SELECT max(\"DatePublished\") FROM \"DeviceLocation\";").Scan(&lastUpdateTimestamp)

	if err != nil {
		log.Err(err).Msg("[Database] Failed to fetch the latest location timestamp")
		return err
	}

	// The query string to insert the data into the database
	queryString := "INSERT INTO \"DeviceLocation\" (\"DeviceID\", \"DatePublished\", \"Description\", \"StatusCode\", \"Payload\") VALUES "

	// Add the value to be inserted
	for i := 0; i < len(fetchedLocations.Results); i++ {
		// Only location the location when there is new value available
		if fetchedLocations.Results[i].DatePublished > lastUpdateTimestamp {
			newLocationNums += 1
			queryString += "('" + fetchedLocations.Results[i].Id + "', " + fmt.Sprintf("%d", fetchedLocations.Results[i].DatePublished) + ", '" + fetchedLocations.Results[i].Description + "', " + fmt.Sprintf("%d", fetchedLocations.Results[i].StatusCode) + ", '" + fetchedLocations.Results[i].Payload + "')"

			// Add a comma to separate each insert data
			if i != len(fetchedLocations.Results)-1 {
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
