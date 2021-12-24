package providers

import (
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/models"
	"time"
)

type FlightProvider interface {
	IsAvailableRoute(srcAirport, destAirport string) bool
	GetOneWayFlights(srcAirport, destAirport string, beginDate, endDate time.Time, numPassengers int, c chan []*models.OneWayFlight)
	CleanUp() //Perform any cleanup needed
}
