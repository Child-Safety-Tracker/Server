package database

type User struct {
	UserID     string `json:"userId"`
	UserName   string `json:"userName"`
	DeviceNums int    `json:"deviceNums"`
}

type Device struct {
	DeviceID   string
	UserID     string
	PrivateKey string
}

type DeviceLocation struct {
	DeviceID      string
	DatePublished int
	Description   string
	StatusCode    int
	Latitude      float32
	Longitude     float32
	Confidence    int
}
