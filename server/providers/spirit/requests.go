package spirit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/constants"
	config2 "github.com/gabriel-flynn/Cheap-Flight-Finder/server/providers/spirit/spirit_config"
	"net/http"
)

func refreshTokenRequest() *http.Request {
	jsonBytes := []byte(`{"credentials":null}`)
	req, _ := http.NewRequest("PUT", constants.RefreshTokenUrl, bytes.NewBuffer(jsonBytes))
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("origin", "https://www.spirit.com")
	req.Header.Set("ocp-apim-subscription-key", "dc6844776fe84b1c8b68affe7deb7916")

	config2.ApiInfo.RLock()
	defer config2.ApiInfo.RUnlock()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config2.ApiInfo.AuthToken))

	return req
}

func generateTokenRequest() *http.Request {
	jsonBytes := []byte(`{"applicationName": "dotRezWeb"}`)
	req, _ := http.NewRequest("POST", constants.GenerateTokenUrl, bytes.NewBuffer(jsonBytes))
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("origin", "https://www.spirit.com")
	req.Header.Set("ocp-apim-subscription-key", "dc6844776fe84b1c8b68affe7deb7916")

	return req
}

func flightRequest(srcAirport, destAirport, startDate, endDate string, numPassengers int, roundTripFlight bool) *http.Request {
	s := NewSearchModel(srcAirport, destAirport, startDate, endDate, numPassengers, roundTripFlight)

	jsonBytes, _ := json.Marshal(s)
	req, _ := http.NewRequest("POST", constants.SearchUrl, bytes.NewBuffer(jsonBytes))
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("origin", "https://www.spirit.com")
	req.Header.Set("ocp-apim-subscription-key", "dc6844776fe84b1c8b68affe7deb7916")

	config2.ApiInfo.RLock()
	defer config2.ApiInfo.RUnlock()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config2.ApiInfo.AuthToken))

	return req
}
