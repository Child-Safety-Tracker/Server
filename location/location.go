package location

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

// Structure of the POST request body
type postRequestBody struct {
	Ids  []string `json:"ids"`
	Days int      `json:"days"`
}

// Structure of each location result
type locationResult struct {
	DatePublished int    `json:"datePublished"`
	Payload       string `json:"payload"`
	Description   string `json:"description"`
	Id            string `json:"id"`
	StatusCode    int    `json:"statusCode"`
}

// Structure of the response body
type postResponseBody struct {
	Results    []locationResult `json:"results"`
	StatusCode string           `json:"statusCode"`
}

func FetchLocation(URL string, id []string, days int) error {
	// Encode the request body
	postBody, err := json.Marshal(&postRequestBody{Ids: id, Days: days})
	if err != nil {
		log.Err(err).Msg("[Location] Error encode the request body")
	}

	// Fire the POST request
	response, err := http.Post(URL, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Err(err).Msg("[Location] Error sending POST request to the Apple Server")
	}

	// Close the response Body when done reading
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Err(err).Msg("[Location] Error reading the response body")
	}

	fmt.Println(string(responseBody))

	postResponseBodyValue := postResponseBody{}
	err = json.Unmarshal(responseBody, &postResponseBodyValue)
	if err != nil {
		log.Err(err).Msg("[Location] Error unmarshalling the response body")
	}

	for index, locationResultValue := range postResponseBodyValue.Results {
		fmt.Println(index, locationResultValue)
	}

	return nil
}
