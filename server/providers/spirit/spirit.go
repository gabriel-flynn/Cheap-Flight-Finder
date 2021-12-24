package spirit

import (
	"encoding/json"
	"fmt"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/models"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/providers"
	"go.uber.org/ratelimit"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type spirit struct {
	rl               ratelimit.Limiter
	quitTokenRefresh chan struct{}
}

var spiritSingleton *spirit
var _ providers.FlightProvider = &spirit{} //Compile-time check that the interface is implemented

var routes map[string]map[string]bool

var doOnce sync.Once

func GetSpiritProvider() *spirit {

	doOnce.Do(func() {
		ticker := time.NewTicker(5 * time.Minute)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					go refreshToken()
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()

		spiritSingleton = &spirit{rl: ratelimit.New(2), quitTokenRefresh: quit}
		_, err := generateToken()
		if err != nil {
			log.Fatal(err)
		}

	})
	return spiritSingleton
}

func (s *spirit) IsAvailableRoute(srcAirport, destAirport string) bool {
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
			for _, route := range data["data"].([]interface{}) {
				key := route.(map[string]interface{})["stationCode"].(string)
				destinations := route.(map[string]interface{})["markets"].([]interface{})
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

func (s *spirit) CleanUp() {
	//Stop refreshing the token
	close(s.quitTokenRefresh)
}

func (s *spirit) GetOneWayFlights(srcAirport, destAirport string, beginDate, endDate time.Time, numPassengers int, c chan []*models.OneWayFlight) {
	//Convert beginDate and endDate to appropriate string format (YYYY-MM-DD)
	localWindowStart := beginDate.Local()
	localWindowEnd := beginDate.Add(time.Hour * 24 * 30).Local() //Spirit allows us to search in upto 31 day windows
	if localWindowEnd.After(endDate) {
		localWindowEnd = endDate
	}

	var wg sync.WaitGroup

	for localWindowStart.Before(endDate) {

		beginDateStr := localWindowStart.Format("2006-01-02")
		endDateStr := localWindowEnd.Format("2006-01-02")

		wg.Add(1)

		go func() {
			defer wg.Done()
			s.rl.Take() //Ensures we don't hit the endpoint with too many TPS
			fmt.Printf("Spirit: %s - %s Querying from %s to %s \n", srcAirport, destAirport, beginDateStr, endDateStr)
			flights, err := getFlightJsonFromAPI(srcAirport, destAirport, beginDateStr, endDateStr, numPassengers, false)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("Spirit: %s - %s Finished Querying from %s to %s \n", srcAirport, destAirport, beginDateStr, endDateStr)
			c <- processFlights(flights, numPassengers)
		}()
		localWindowStart = localWindowEnd.Add(time.Hour * 24)
		localWindowEnd = localWindowStart.Add(time.Hour * 24 * 30)
		if localWindowEnd.After(endDate) {
			localWindowEnd = endDate
		}
	}

	wg.Wait()
}

func processFlights(trips []interface{}, numPassengers int) []*models.OneWayFlight {
	//trips is an array with two items in it -> item #1 contains all flight information from srcAirport to destAirport and item #2 contains all flight information from destAirport to srcAirport
	srcTrip := trips[0]
	srcFlights := srcTrip.(map[string]interface{})["journeysAvailable"].([]interface{})
	srcAiport := srcTrip.(map[string]interface{})["origin"].(string)
	destAirport := srcTrip.(map[string]interface{})["destination"].(string)

	var flights []*models.OneWayFlight

	for _, flight := range srcFlights {
		srcFare, fareKey := getFare(flight.(map[string]interface{})["fares"].(map[string]interface{}), FlightRequirements.IsSaversClub, numPassengers) //Get the standard fare, not the saver's club fare

		// If srcFare is nil then that means none of the fares met our requirements (ex: number of seats available)
		if srcFare == nil {
			continue
		}

		//Get the fare price and fare key
		passengerFares := srcFare.(map[string]interface{})["details"].(map[string]interface{})["passengerFares"].([]interface{})

		//Get the cheapest passenger fare with the number of seats requirements we need
		var cheapestFare map[string]interface{} = nil
		for _, passengerFare := range passengerFares {

			if cheapestFare == nil || passengerFare.(map[string]interface{})["fareAmount"].(float64) < cheapestFare["fareAmount"].(float64) {
				cheapestFare = passengerFare.(map[string]interface{})
			}

		}

		departureTime, err := time.Parse("2006-01-02T15:04:05", flight.(map[string]interface{})["designator"].(map[string]interface{})["departure"].(string))
		if err != nil {
			log.Println(err)
			continue
		}

		arrivalTime, err := time.Parse("2006-01-02T15:04:05", flight.(map[string]interface{})["designator"].(map[string]interface{})["arrival"].(string))
		if err != nil {
			log.Println(err)
			continue
		}

		//Convert duration string to
		duration := flight.(map[string]interface{})["duration"].(string)
		durationSlice := strings.Split(duration, ":")
		if len(durationSlice) != 3 {
			log.Fatalf("Duration was not formatted as expected: %s", duration)
		}
		flightDuration, err := time.ParseDuration(fmt.Sprintf("%sh%sm%ss", durationSlice[0], durationSlice[1], durationSlice[2]))
		if err != nil {
			log.Println(err)
			continue
		}

		flights = append(flights, &models.OneWayFlight{
			FlightKey:       flight.(map[string]interface{})["journeyKey"].(string),
			Airline:         "spirit",
			SrcAirport:      srcAiport,
			DestAirport:     destAirport,
			FareKey:         fareKey,
			FarePrice:       float32(cheapestFare["fareAmount"].(float64)),
			FlightDeparture: departureTime,
			FlightArrival:   arrivalTime,
			FlightDuration:  flightDuration,
		})
	}

	return flights
}

func getFare(fares map[string]interface{}, isSaversClub bool, numPassengers int) (interface{}, string) {
	for fareKey := range fares {
		fareDetails := fares[fareKey].(map[string]interface{})["details"].(map[string]interface{})
		if fareDetails["isClubFare"].(bool) == isSaversClub && int(fares[fareKey].(map[string]interface{})["availableCount"].(float64)) >= numPassengers {
			return fares[fareKey], fareKey
		}
	}
	log.Fatalf("An unexcepted error occured with the following fare: %v", fares)
	return nil, ""
}
