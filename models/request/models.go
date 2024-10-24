package request

// LocationRequest - Structure of the POST request body
type LocationRequest struct {
	PrivateKey string   `json:"privateKey"`
	Ids        []string `json:"ids"`
	Days       int      `json:"days"`
}
