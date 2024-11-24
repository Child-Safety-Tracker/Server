package main

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"server/database"
	"server/models/request"
	responseModels "server/models/response"
	"time"
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

func main() {

	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("[Server] Error loading .env file")
		return
	}

	// Logger init with pretty format and timestamp enabled
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
	// Set the logger to use globally
	log.Logger = logger

	// Connect to the database
	_ = database.DatabaseConnect()
}
