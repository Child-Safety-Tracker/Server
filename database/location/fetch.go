package location

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"server/misc/decrypt"
	locationModels "server/models/location"
	"server/models/request"
	responseModels "server/models/response"
	"slices"
	"strconv"
)

func AppleServerLocations(ids []string, days int) (responseModels.LocationResponse, error) {

	// Apple server wrapper URL
	URL := os.Getenv("APPLE_SERVER_WRAPPER_URL")

	// Encode the request body
	postBody, err := json.Marshal(&request.LocationRequest{Ids: ids, Days: days})
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

func FetchLocation(database *pgxpool.Pool, ids []string, privateKeys []string) ([]locationModels.DecryptedLocationResult, error) {
	var decryptedLocationResultValue []locationModels.DecryptedLocationResult

	fetchedLocations, err := AppleServerLocations(ids, 1)

	if err != nil {
		log.Error().Msg("[Location] Error fetching location from Apple server")
		return []locationModels.DecryptedLocationResult{}, err
	}

	// There is no location published
	if len(fetchedLocations.Results) == 0 {
		log.Warn().Msg("[Location] Empty fetched location")
		return []locationModels.DecryptedLocationResult{}, err
	}

	// Split results by device
	// 5 indicates the maximum devices
	var splitLocations [5][]locationModels.LocationResult
	for _, element := range fetchedLocations.Results {
		for idIndex, id := range ids {
			if element.Id == id {
				splitLocations[idIndex] = append(splitLocations[idIndex], element)
			}
		}
	}

	// Get the latest location from split results
	var latestLocations [5]locationModels.LocationResult
	for index, elements := range splitLocations {
		if len(splitLocations[index]) > 1 {
			latestLocations[index] = slices.MaxFunc(elements, func(i, j locationModels.LocationResult) int {
				return cmp.Compare(i.DatePublished, j.DatePublished)
			})
		} else if len(splitLocations[index]) == 1 {
			latestLocations[index] = elements[0]
		}
	}

	// Check if the latestLocation contains all the requested IDs
	for index, element := range ids {
		// If the latest location for that index is empty
		if latestLocations[index].Id != element {
			fmt.Println(element)
			var queriedLocation locationModels.LocationResult
			var locationId int

			err := database.QueryRow(context.Background(), "SELECT * FROM \"DeviceLocation\" WHERE \"DatePublished\"=(SELECT MAX(\"DatePublished\") FROM \"DeviceLocation\" WHERE \"DeviceID\"=$1);", element).Scan(&locationId, &queriedLocation.Id, &queriedLocation.DatePublished, &queriedLocation.Description, &queriedLocation.StatusCode, &queriedLocation.Payload)

			if err != nil {
				log.Err(err).Msg("[Location] Error getting latest location from the database")
				return []locationModels.DecryptedLocationResult{}, err
			}

			// Update to the latest location from the database
			latestLocations[index] = queriedLocation
		}
	}

	// Decrypt it
	for index, element := range latestLocations {
		// temporary variable to handle error
		if element != (locationModels.LocationResult{}) {
			fmt.Printf("%+v\n", element)
			tempDecryptedLocation, err := decrypt.DecryptLocation(element, privateKeys[index])

			if err != nil {
				log.Error().Msg("[Location] Error decrypting the location payload at index: " + strconv.Itoa(index))
			}

			decryptedLocationResultValue = append(decryptedLocationResultValue, tempDecryptedLocation)
		}
	}

	return decryptedLocationResultValue, nil

}
