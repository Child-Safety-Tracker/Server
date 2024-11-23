package update

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"server/models/request"
	responseModels "server/models/response"
)

func AppleServerLocations(URL string, id []string, days int) (responseModels.LocationResponse, error) {
	// Encode the request body
	postBody, err := json.Marshal(&request.LocationRequest{Ids: id, Days: days})
	if err != nil {
		log.Err(err).Msg("[Location] Error encode the request body")
		return responseModels.LocationResponse{}, err
	}

	// Fire the POST request
	response, err := http.Post(URL, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Err(err).Msg("[Location] Error sending POST request to the Apple Server")
		return responseModels.LocationResponse{}, err
	}

	// Close the response Body when done reading
	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Err(err).Msg("[Location] Error reading the response body")
		return responseModels.LocationResponse{}, err
	}

	// Unmarshall the json response body
	postResponseBodyValue := responseModels.LocationResponse{}
	err = json.Unmarshal(responseBody, &postResponseBodyValue)
	if err != nil {
		log.Err(err).Msg("[Location] Error unmarshalling the response body")
		return responseModels.LocationResponse{}, err
	}

	return postResponseBodyValue, nil
}

func DatabaseUpdateLocation(database *pgx.Conn, ids []string) error {
	newLocationNums := 0 // Indicate the number of new locations
	// Fetch the locations from Apple server
	response, err := AppleServerLocations(os.Getenv("APPLE_SERVER_WRAPPER_URL"), ids, 7)

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
		// Only update the location when there is new value available
		if response.Results[i].DatePublished > lastUpdateTimestamp {
			newLocationNums += 1
			queryString += "('" + response.Results[i].Id + "', " + fmt.Sprintf("%d", response.Results[i].DatePublished) + ", '" + response.Results[i].Description + "', " + fmt.Sprintf("%d", response.Results[i].StatusCode) + ", '" + response.Results[i].Payload + "')"

			// Add a comma to separate each insert data
			if i != len(response.Results)-1 {
				queryString += ","
			}
		}
		// Trim the last trailing comma
		if newLocationNums != 0 {
			queryString = queryString[:len(queryString)-1]
		}
	}

	if newLocationNums == 0 {
		// No update when there is no new location
		queryString = ""
		log.Warn().Msg("[Database] No new location to be updated")
	}

	// Execute the insert command
	_, err = database.Exec(context.Background(), queryString)

	if err != nil {
		log.Err(err).Msg("[Database] Error executing the update location query command")
		return err
	}

	log.Info().Msg("[Database] Successfully updated the location")
	return nil
}
