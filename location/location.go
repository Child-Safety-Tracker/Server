package location

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

type requestBody struct {
	Ids  []string `json:"ids"`
	Days int      `json:"days"`
}

func FetchLocation(URL string, id []string, days int) error {
	// Encode the request body
	postBody, err := json.Marshal(&requestBody{Ids: id, Days: days})
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

	return nil
}
