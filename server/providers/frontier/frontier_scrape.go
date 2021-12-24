package frontier

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/rpc"
	"log"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`FlightData = '({.*})';`)

func getFlights(url string) []interface{} {

	client := rpc.GetClient()
	req := rpc.PageSourceRequest{
		Url: url,
	}
	resp, err := client.GetPageSource(context.Background(), &req)
	if err != nil {
		log.Fatal(err)
	}

	flightMatch := re.FindStringSubmatch(resp.PageSource)

	//TODO: Use the avast-go retry to only retry a certain number of times
	if len(flightMatch) == 0 {
		return getFlights(url)
	}
	flightJson := strings.ReplaceAll(flightMatch[1], "&quot;", "\"")

	//Convert json to map
	var flightMap map[string]interface{}

	err = json.Unmarshal([]byte(flightJson), &flightMap)
	if err != nil {
		fmt.Println("Unable to unmarshal json: ", err)
	}
	return flightMap["journeys"].([]interface{})
}

func getFlightsLambda(url string) []interface{} {
	payload := fmt.Sprintf("{\"url\": \"%s\"}", url)
	flightJson := rpc.GetPageSourceLambda([]byte(payload))

	//Convert json to map
	var flightMap map[string]interface{}

	err := json.Unmarshal([]byte(flightJson), &flightMap)
	if err != nil {
		fmt.Println("Unable to unmarshal json: ", err)
	}

	//TODO: Use avast-go to retry
	journeys, ok := flightMap["journeys"].([]interface{})
	if !ok {
		if flightMap["errorMessage"] != nil {
			return getFlightsLambda(url)
		} else {
			log.Fatal(flightJson)
		}
	}
	return journeys
}
