package southwest

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/models"
	"go.uber.org/ratelimit"
	"golang.org/x/sync/semaphore"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type southwest struct {
	rl  ratelimit.Limiter
	sem *semaphore.Weighted
}

var southwestSingleton *southwest

var doOnce sync.Once

func GetSouthwestProvider() *southwest {

	//TODO: revisit if using reusing auth headers works well for southwest. In the latest attempt we would use one set of auth headers per request which bottlenecks the our throughput quite heavily
	doOnce.Do(func() {
		southwestSingleton = &southwest{rl: ratelimit.New(2), sem: semaphore.NewWeighted(5)}
	})
	return southwestSingleton
}

func (s *southwest) CleanUp() {
	//Stop refreshing the token
}

func (s *southwest) GetOneWayFlights(srcAirport, destAirport string, beginDate, endDate time.Time, numPassengers int, c chan []*models.OneWayFlight) {
	//Convert beginDate to appropriate string format (YYYY-MM-DD)
	localWindowStart := beginDate.Local() //Southwest flights have to be queried day by day AFAIK (Can't batch request like with Spirit)

	var wg sync.WaitGroup

	for localWindowStart.Before(endDate) {

		beginDateStr := localWindowStart.Format("2006-01-02")

		wg.Add(1)

		go func() {
			defer wg.Done()
			s.rl.Take() //Ensures we don't hit the endpoint with too many TPS
			fmt.Printf("Southwest: %s - %s Querying from %s to %s \n", srcAirport, destAirport, beginDateStr, beginDateStr)

			//time.Sleep(time.Minute * 3)
			var flights []interface{}
			var err error
			err = retry.Do(
				func() error {
					s.sem.Acquire(context.TODO(), 1)
					refreshAuthHeaders()
					fmt.Printf("Trying to get flight data for date %s", beginDateStr)
					flights, err = getFlightJsonFromAPI(srcAirport, destAirport, beginDateStr, beginDateStr, numPassengers, false)
					s.sem.Release(1)
					if err != nil {
						fmt.Println("ERROR!", err)
						return err
					}
					fmt.Println("SUCCESS!")
					return nil
				},
				retry.Attempts(5),
				retry.Delay(time.Minute))

			fmt.Printf("Southwest: %s - %s Finished Querying from %s to %s \n", srcAirport, destAirport, beginDateStr, beginDateStr)
			c <- processFlights(flights, beginDateStr)
		}()
		localWindowStart = localWindowStart.Add(time.Hour * 24)
	}

	wg.Wait()
}

func processFlights(flightsMap []interface{}, date string) []*models.OneWayFlight {

	var flights []*models.OneWayFlight

	for _, flight := range flightsMap {
		// No fares available
		if flight.(map[string]interface{})["fares"] == nil {
			continue
		}
		srcFare := getLowestFare(flight.(map[string]interface{})["fares"].([]interface{}))

		// If srcFare is nil then that means none of the fares met our requirements (ex: number of seats available)
		if srcFare == nil {
			continue
		}

		//Get the fare price
		passengerFarePrice, _ := strconv.ParseFloat(srcFare.(map[string]interface{})["price"].(map[string]interface{})["amount"].(string), 32)

		departureTime, err := time.Parse("2006-01-02T15:04", fmt.Sprintf("%sT%s", date, flight.(map[string]interface{})["departureTime"].(string)))
		if err != nil {
			log.Println(err)
			continue
		}

		arrivalTime, err := time.Parse("2006-01-02T15:04", fmt.Sprintf("%sT%s", date, flight.(map[string]interface{})["arrivalTime"].(string)))
		if err != nil {
			log.Println(err)
			continue
		}

		//Convert duration string
		duration := flight.(map[string]interface{})["duration"].(string)
		flightDuration, err := time.ParseDuration(strings.ReplaceAll(duration, " ", ""))
		if err != nil {
			log.Println(err)
			continue
		}

		flightKey := flight.(map[string]interface{})["_meta"].(map[string]interface{})["cardId"].(string)
		flightSlice := strings.Split(flightKey, ":")
		srcAirport := flightSlice[0]
		destAirport := flightSlice[1]

		flights = append(flights, &models.OneWayFlight{
			FlightKey:       flight.(map[string]interface{})["_meta"].(map[string]interface{})["cardId"].(string),
			Airline:         "southwest",
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

func getLowestFare(fares []interface{}) interface{} {
	var fare map[string]interface{}
	for fareKey := range fares {

		priceMap := fares[fareKey].(map[string]interface{})["price"]
		if priceMap == nil {
			continue
		}

		farePrice, err := strconv.ParseFloat(priceMap.(map[string]interface{})["amount"].(string), 32)
		if err != nil {
			fmt.Println("Couldn't convert farePrice to float.", err)
		}

		var lowestFarePrice float64
		if fare != nil {
			lowestFarePrice, err = strconv.ParseFloat(fare["price"].(map[string]interface{})["amount"].(string), 32)
		}

		if err != nil {
			fmt.Println("Couldn't convert lowestFarePrice to float.", err)
		}

		if fares[fareKey].(map[string]interface{})["price"] != nil && fare == nil || farePrice < lowestFarePrice {
			fare = fares[fareKey].(map[string]interface{})
		}
	}
	return fare
}
