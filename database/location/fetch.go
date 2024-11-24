package location

import (
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"server/misc/decrypt"
	locationModels "server/models/location"
	"server/models/request"
	responseModels "server/models/response"
)

func AppleServerLocations(id []string, days int) (responseModels.LocationResponse, error) {

	// Apple server wrapper URL
	URL := os.Getenv("APPLE_SERVER_WRAPPER_URL")

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

func FetchLocation(ids []string, privateKey string) (locationModels.DecryptedLocationResult, error) {
	decryptedLocationResultValue := locationModels.DecryptedLocationResult{}

	fetchedLocations, err := AppleServerLocations(ids, 7)

	if err != nil {
		log.Error().Msg("[Location] Error fetching location from Apple server")
		return locationModels.DecryptedLocationResult{}, err
	}

	// There is no location published
	if len(fetchedLocations.Results) == 0 {
		log.Warn().Msg("[Location] Empty fetched location")
		return locationModels.DecryptedLocationResult{}, err
	}

	decryptedLocationResultValue, err = decrypt.DecryptLocation(fetchedLocations.Results[0], privateKey)

	if err != nil {
		log.Error().Msg("[Location] Error decrypting the location payload")
		return locationModels.DecryptedLocationResult{}, err
	}

	return decryptedLocationResultValue, nil

}
