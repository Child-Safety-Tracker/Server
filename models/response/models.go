package response

import "server/models/location"

// LocationResponse - Structure of the response body
type LocationResponse struct {
	Results    []location.LocationResult `json:"results"`
	StatusCode string                    `json:"statusCode"`
}
