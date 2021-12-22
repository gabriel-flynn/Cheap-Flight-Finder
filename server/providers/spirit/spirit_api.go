package spirit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func generateToken() (string, error) {
	client := &http.Client{}
	res, err := client.Do(generateTokenRequest())
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var x map[string]interface{}
		json.Unmarshal(bodyBytes, &x)

		token := x["data"].(map[string]interface{})["token"].(string)
		if ApiInfo == nil {
			ApiInfo = &ApiConfig{AuthToken: token}
		}

		return token, nil

	}

	bodyBufr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("an unexpected error occurred when trying to refresh the token")
	}

	return "", errors.New(fmt.Sprintf("an unexpected error occurred when trying to refresh the token: %s", string(bodyBufr)))
}

func refreshToken() {
	client := &http.Client{}
	res, err := client.Do(refreshTokenRequest())
	if err != nil {
		log.Fatal(err)
	}

	res.Body.Close()
}

func getFlightJsonFromAPI(srcAirport, destAirport, startDate, endDate string, numPassengers int, roundTripFlight bool) ([]interface{}, error) {
	client := &http.Client{}
	res, err := client.Do(flightRequest(srcAirport, destAirport, startDate, endDate, numPassengers, roundTripFlight))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}


		var x map[string]interface{}
		err = json.Unmarshal(bodyBytes, &x)
		if err != nil {
			log.Fatal(err)
		}

		return x["data"].(map[string]interface{})["trips"].([]interface{}), nil
	}


	bodyBufr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("an unexpected error occurred when trying to refresh the token")
	}

	return nil, errors.New(fmt.Sprintf("an unexpected error occurred when trying to retrieve the flights: %s", string(bodyBufr)))
}
