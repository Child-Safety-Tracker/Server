package request

// LocationRequest - Structure of the POST request body
type LocationRequest struct {
	PrivateKeys []string `json:"privateKeys"`
	Ids         []string `json:"ids"`
	Days        int      `json:"days"`
}

type DeviceStatusEditRequest struct {
	DeviceId string `json:"deviceId"`
	Enabled  bool   `json:"enabled"`
}
