package spirit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func refreshTokenRequest() *http.Request {
	jsonBytes := []byte(`{"credentials":null}`)
	req, _ := http.NewRequest("PUT", RefreshTokenUrl, bytes.NewBuffer(jsonBytes))
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("origin", "https://www.spirit.com")
	req.Header.Set("ocp-apim-subscription-key", "dc6844776fe84b1c8b68affe7deb7916")

	ApiInfo.RLock()
	defer ApiInfo.RUnlock()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ApiInfo.AuthToken))

	return req
}

func generateTokenRequest() *http.Request {
	jsonBytes := []byte(`{"applicationName": "dotRezWeb"}`)
	req, _ := http.NewRequest("POST", GenerateTokenUrl, bytes.NewBuffer(jsonBytes))
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("origin", "https://www.spirit.com")
	req.Header.Set("ocp-apim-subscription-key", "dc6844776fe84b1c8b68affe7deb7916")

	return req
}

func flightRequest(srcAirport, destAirport, startDate, endDate string, numPassengers int, roundTripFlight bool) *http.Request {
	s := NewSearchModel(srcAirport, destAirport, startDate, endDate, numPassengers, roundTripFlight)

	jsonBytes, _ := json.Marshal(s)
	req, _ := http.NewRequest("POST", SearchUrl, bytes.NewBuffer(jsonBytes))
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("origin", "https://www.spirit.com")
	req.Header.Set("ocp-apim-subscription-key", "dc6844776fe84b1c8b68affe7deb7916")

	ApiInfo.RLock()
	defer ApiInfo.RUnlock()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ApiInfo.AuthToken))

	return req
}
