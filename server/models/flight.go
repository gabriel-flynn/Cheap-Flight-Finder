package models

import (
	"fmt"
	"time"
)

type OneWayFlight struct {
	FlightKey       string
	Airline         string
	SrcAirport      string
	DestAirport     string
	FareKey         string
	FarePrice       float32
	FlightDeparture time.Time
	FlightArrival   time.Time
	FlightDuration  time.Duration
}

func (f *OneWayFlight) String() string {
	return fmt.Sprintf("Flight[FlightKey: %s, Airline: %s, SrcAiport: %s, DestAirport: %s, FareKey: %s, FarePrice: %f, FlightDeparture: %s, FlightArrival: %s, FlightDuration: %s]",
		f.FlightKey, f.Airline, f.SrcAirport, f.DestAirport, f.FareKey, f.FarePrice, f.FlightDeparture, f.FlightArrival, f.FlightDuration)
}

func (f *OneWayFlight) GetOneWayFlightRecord() []string {
	farePrice := fmt.Sprintf("%.2f", f.FarePrice)
	depTime := f.FlightDeparture.Format("2006-01-02 15:04:05")
	arrTime := f.FlightArrival.Format("2006-01-02 15:04:05")
	return []string{f.FlightKey, f.Airline, f.SrcAirport, f.DestAirport, f.FareKey, farePrice,
		depTime, arrTime, f.FlightDuration.String()}
}

func BatchGetOneWayFlightRecord(flights []*OneWayFlight) [][]string {
	var flightRecords [][]string
	for _, flight := range flights {
		flightRecords = append(flightRecords, flight.GetOneWayFlightRecord())
	}
	return flightRecords
}

func GetOneWayFlightHeading() []string {
	return []string{"FlightKey", "Airline", "SrcAirport", "DestAirport", "FareKey", "FarePrice", "FlightDeparture", "FlightArrival", "FlightDuration"}
}

type RoundTripFlight struct {
	FareKey             string
	Airline             string
	SrcAirport          string
	DestAirport         string
	FarePrice           float32
	SrcFlightDeparture  time.Time
	SrcFlightArrival    time.Time
	SrcFlightDuration   time.Duration
	DestFlightDeparture time.Time
	DestFlightArrival   time.Time
	DestFlightDuration  time.Duration
}

func (f *RoundTripFlight) String() string {
	return fmt.Sprintf("Flight[FareKey: %s, Airline: %s, SrcAiport: %s, DestAirport: %s, FarePrice: %f, SrcFlightDeparture: %s, SrcFlightArrival: %s, SrcFlightDuration: %s, DestFlightDeparture: %s, DestFlightArrival: %s, DestFlightDuration: %s ]",
		f.FareKey, f.Airline, f.SrcAirport, f.DestAirport, f.FarePrice, f.SrcFlightDeparture, f.SrcFlightArrival, f.SrcFlightDuration, f.DestFlightDeparture, f.DestFlightArrival, f.DestFlightDuration)
}
