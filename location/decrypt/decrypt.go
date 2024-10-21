package decrypt

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os/exec"
	"server/location/models"
)

func DecryptLocation(locationResultValue models.LocationResult) (models.DecryptedLocationResult, error) {

	// Exec the decryption script on the payload
	decodeOutput, err := exec.Command("bash", "-c", "python3 location/decrypt/decrypt.py "+locationResultValue.Payload).Output()
	if err != nil {
		log.Fatal().Err(err).Msg("[Location] Failed to decrypt payload")
		return models.DecryptedLocationResult{}, err
	}

	decryptedLocation := models.DecryptedLocation{}
	err = json.Unmarshal(decodeOutput, &decryptedLocation)
	if err != nil {
		log.Err(err).Msg("[Location] Error unmarshalling the decrypted location")
		return models.DecryptedLocationResult{}, err
	}

	return models.DecryptedLocationResult{
		DatePublished: locationResultValue.DatePublished,
		Payload:       decryptedLocation,
		Description:   locationResultValue.Description,
		Id:            locationResultValue.Id,
		StatusCode:    locationResultValue.StatusCode,
	}, nil
}
