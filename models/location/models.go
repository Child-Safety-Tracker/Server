package location

// LocationResult - Structure of each location result
type LocationResult struct {
	DatePublished int    `json:"datePublished"`
	Payload       string `json:"payload"`
	Description   string `json:"description"`
	Id            string `json:"id"`
	StatusCode    int    `json:"statusCode"`
}

// DecryptedLocation - Structure of decrypted location and the general result
type DecryptedLocation struct {
	Longitude  float32 `json:"longitude"`
	Latitude   float32 `json:"latitude"`
	Confidence int     `json:"confidence"`
	Timestamp  int     `json:"timestamp"`
}
type DecryptedLocationResult struct {
	DatePublished int               `json:"datePublished"`
	Payload       DecryptedLocation `json:"payload"`
	Description   string            `json:"description"`
	Id            string            `json:"id"`
	StatusCode    int               `json:"statusCode"`
}
