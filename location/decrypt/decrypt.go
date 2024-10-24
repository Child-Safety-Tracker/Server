package decrypt

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os/exec"
	"server/models/location"
)

func DecryptLocation(locationResultValue location.LocationResult, privateKey string) (location.DecryptedLocationResult, error) {

	// Exec the decryption script on the payload
	decodeOutput, err := exec.Command("bash", "-c", "python3 location/decrypt/decrypt.py "+privateKey+" "+locationResultValue.Payload).Output()
	if err != nil {
		log.Err(err).Msg("[Location] Failed to decrypt payload")
		return location.DecryptedLocationResult{}, err
	}

	decryptedLocation := location.DecryptedLocation{}
	err = json.Unmarshal(decodeOutput, &decryptedLocation)
	if err != nil {
		log.Err(err).Msg("[Location] Error unmarshalling the decrypted location")
		return location.DecryptedLocationResult{}, err
	}

	return location.DecryptedLocationResult{
		DatePublished: locationResultValue.DatePublished,
		Payload:       decryptedLocation,
		Description:   locationResultValue.Description,
		Id:            locationResultValue.Id,
		StatusCode:    locationResultValue.StatusCode,
	}, nil
}
