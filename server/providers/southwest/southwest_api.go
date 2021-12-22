package southwest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func getFlightJsonFromAPI(srcAirport, destAirport, startDate, endDate string, numPassengers int, roundTripFlight bool) ([]interface{}, error) {
	client := &http.Client{}

	var res *http.Response
	var err error
	if !roundTripFlight {
		res, err = client.Do(oneWayFlightRequest(srcAirport, destAirport, startDate, numPassengers, roundTripFlight))
	}
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

		return x["flightShoppingPage"].(map[string]interface{})["outboundPage"].(map[string]interface{})["cards"].([]interface{}), nil
	}


	bodyBufr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("an unexpected error occurred when trying to refresh the token")
	}

	return nil, errors.New(fmt.Sprintf("an unexpected error occurred when trying to retrieve the flights: %s", string(bodyBufr)))
}
