package models

type User struct {
	UserID     string
	Username   string
	DeviceNums int
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
