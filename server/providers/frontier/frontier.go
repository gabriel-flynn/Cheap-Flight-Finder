package frontier

import (
	"encoding/json"
	"fmt"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/models"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/providers"
	"go.uber.org/ratelimit"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type frontier struct {
	rl ratelimit.Limiter
}

var frontierSingleton *frontier
var _ providers.FlightProvider = &frontier{} //Compile-time check that the interface is implemented

var routes map[string]map[string]bool

var doOnce sync.Once

func GetFrontierProvider() *frontier {

	doOnce.Do(func() {
		frontierSingleton = &frontier{rl: ratelimit.New(5)}
	})
	return frontierSingleton
}

func (s *frontier) IsAvailableRoute(srcAirport, destAirport string) bool {
	if routes == nil {
		client := &http.Client{}
		res, err := client.Do(getRoutesRequest())
		if err != nil {
			log.Fatal("Unable to successfully get Frontier routes", err)
		}

		defer res.Body.Close()
		if res.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(res.Body)
			if err != nil {
				log.Fatal("Unable to successfully get Frontier routes body", err)
			}

			var data map[string]interface{}
			err = json.Unmarshal(bodyBytes, &data)
			if err != nil {
				log.Fatal("Unable to successfully unmarshal Frontier routes JSON", err)
			}

			routes = make(map[string]map[string]bool)
			for _, route := range data["markets"].([]interface{}) {
				key := route.(map[string]interface{})["fromStation"].(string)
				destinations := route.(map[string]interface{})["toStations"].([]interface{})
				set := make(map[string]bool)
				for _, dest := range destinations {
					set[dest.(string)] = true
				}
				routes[key] = set
			}
		}
	}

	if val, ok := routes[srcAirport]; ok {
		if _, ok := val[destAirport]; ok {
			return true
		}
	}

	return false
}

func (s *frontier) CleanUp() {
}

func (s *frontier) GetOneWayFlights(srcAirport, destAirport string, beginDate, endDate time.Time, numPassengers int, c chan []*models.OneWayFlight) {
	//Convert beginDate to appropriate string format (YYYY-MM-DD)
	localWindowStart := beginDate.Local() //Southwest flights have to be queried day by day AFAIK (Can't batch request like with Spirit)

	var wg sync.WaitGroup

	for localWindowStart.Before(endDate) {

		beginDateStr := localWindowStart.Format("2006-01-02")

		wg.Add(1)

		go func(localWindowStart time.Time) {
			defer wg.Done()
			s.rl.Take() //Ensures we don't hit the endpoint with too many TPS
			fmt.Printf("Frontier: %s - %s Querying from %s to %s \n", srcAirport, destAirport, beginDateStr, beginDateStr)

			var flights []interface{}

			//Create url encoding
			params := url.Values{}
			params.Add("o1", srcAirport)
			params.Add("d1", destAirport)
			params.Add("dd1", localWindowStart.Format("Jan 2, 2006"))
			params.Add("ADT", strconv.FormatInt(int64(numPassengers), 10))
			params.Add("mon", "true")
			params.Add("promo", "") //TODO: Get promo from config

			fmt.Printf("Trying to get flight data for date %s\n", beginDateStr)
			flights = getFlightsLambda(FrontierBookingUrl + params.Encode())

			fmt.Printf("Frontier: %s - %s Finished Querying from %s to %s \n", srcAirport, destAirport, beginDateStr, beginDateStr)
			c <- processFlights(flights, beginDateStr)
		}(localWindowStart)
		localWindowStart = localWindowStart.Add(time.Hour * 24)
	}

	wg.Wait()
}

func processFlights(journeysMap []interface{}, date string) []*models.OneWayFlight {

	var flights []*models.OneWayFlight

	//There is always only one item for one-way flights (hence index 0)
	flightsMap := journeysMap[0].(map[string]interface{})["flights"].([]interface{})
	for _, flight := range flightsMap {

		flightKey := flight.(map[string]interface{})["standardFareKey"].(string)

		//Get the fare price
		passengerFarePrice := flight.(map[string]interface{})["standardFare"].(float64)

		if flightKey == "" {
			continue
		}

		var otpData []interface{}
		json.Unmarshal([]byte(flight.(map[string]interface{})["otpData"].(string)), &otpData)

		departureTime, err := time.Parse("2006-01-02 15:04", otpData[0].(map[string]interface{})["departureUTCDateTime"].(string))
		if err != nil {
			log.Println(err)
			continue
		}

		//TODO: Get the correct time zone -> we're just putting central right now but that isn't accurate
		centralTime, _ := time.LoadLocation("America/Chicago")
		arrivalTime, err := time.ParseInLocation("Monday, January _2, 2006 3:04PM", otpData[0].(map[string]interface{})["arrivalDate"].(string)+" "+otpData[0].(map[string]interface{})["arrivalTime"].(string), centralTime)
		if err != nil {
			log.Println(err)
			continue
		}

		//Convert duration string
		duration := flight.(map[string]interface{})["duration"].(string)
		duration = strings.ReplaceAll(duration, "hrs", "h")
		duration = strings.ReplaceAll(duration, "min", "m")
		flightDuration, err := time.ParseDuration(strings.ReplaceAll(duration, " ", ""))
		if err != nil {
			log.Println(err)
			continue
		}

		srcAirport := journeysMap[0].(map[string]interface{})["departureStation"].(string)
		destAirport := journeysMap[0].(map[string]interface{})["arrivalStation"].(string)

		flights = append(flights, &models.OneWayFlight{
			FlightKey:       flightKey,
			Airline:         "frontier",
			SrcAirport:      srcAirport,
			DestAirport:     destAirport,
			FareKey:         "",
			FarePrice:       float32(passengerFarePrice),
			FlightDeparture: departureTime,
			FlightArrival:   arrivalTime,
			FlightDuration:  flightDuration,
		})
	}

	return flights
}
