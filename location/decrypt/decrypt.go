package decrypt

import (
	"encoding/base64"
	"github.com/rs/zerolog/log"
	"os/exec"
)

func DecryptLocation(payload string) {

	_, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		log.Fatal().Err(err).Msg("[Location] Failed to decode payload")
	}

	decodeOutput, err := exec.Command("bash", "-c", "python3 location/decrypt/decrypt.py "+payload).Output()
	if err != nil {
		log.Fatal().Err(err).Msg("[Location failed to decrypt payload]")
	}
	println(string(decodeOutput))

}
