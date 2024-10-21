package location

import (
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"server/location/models"
)

func FetchLocation(URL string, id []string, days int) (models.PostResponseBody, error) {
	// Encode the request body
	postBody, err := json.Marshal(&models.PostRequestBody{Ids: id, Days: days})
	if err != nil {
		log.Err(err).Msg("[Location] Error encode the request body")
		return models.PostResponseBody{}, err
	}

	// Fire the POST request
	response, err := http.Post(URL, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Err(err).Msg("[Location] Error sending POST request to the Apple Server")
		return models.PostResponseBody{}, err
	}

	// Close the response Body when done reading
	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Err(err).Msg("[Location] Error reading the response body")
		return models.PostResponseBody{}, err
	}

	// Unmarshall the json response body
	postResponseBodyValue := models.PostResponseBody{}
	err = json.Unmarshal(responseBody, &postResponseBodyValue)
	if err != nil {
		log.Err(err).Msg("[Location] Error unmarshalling the response body")
		return models.PostResponseBody{}, err
	}

	return postResponseBodyValue, nil
}
