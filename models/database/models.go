package database

type User struct {
	UserID     string `json:"userId"`
	UserName   string `json:"userName"`
	DeviceNums int    `json:"deviceNums"`
}

type Device struct {
	DeviceID   string `json:"deviceId"`
	UserID     string `json:"userId"`
	PrivateKey string `json:"privateKey"`
	Enabled    bool   `json:"enabled"`
}

type DeviceLocation struct {
	LocationID    string
	DeviceID      string
	DatePublished int
	Description   string
	StatusCode    int
	Payload       string
}
