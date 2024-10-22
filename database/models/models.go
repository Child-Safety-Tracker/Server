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
