package coordinator

import (
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/config"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/models"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/providers"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/providers/frontier"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/providers/spirit"
	"sync"
	"time"
)

type flightProviderCoordinator struct {
	aggregate chan []*models.OneWayFlight
}

var coordinatorSingleton *flightProviderCoordinator

var doOnce sync.Once

func NewFlightProviderCoordinator(channel chan []*models.OneWayFlight) *flightProviderCoordinator {
	doOnce.Do(func() {
		coordinatorSingleton = &flightProviderCoordinator{channel}
	})
	return coordinatorSingleton
}

func (f *flightProviderCoordinator) Start() {
	var chans []chan []*models.OneWayFlight
	flightProviders := []providers.FlightProvider{spirit.GetSpiritProvider(), frontier.GetFrontierProvider()}

	for indx, provider := range flightProviders {
		c := make(chan []*models.OneWayFlight)
		chans = append(chans, c)
		flightSearchConfig := config.GetConfig()
		var wg sync.WaitGroup
		for _, f := range flightSearchConfig.FlightParameters {
			if provider.IsAvailableRoute(f.SrcAirport, f.DestAirport) {
				wg.Add(1)
				go func(f config.FlightSearch, wg *sync.WaitGroup, indx int, provider providers.FlightProvider) {
					defer wg.Done()
					provider.GetOneWayFlights(f.SrcAirport, f.DestAirport, f.StartWindow.Shift(time.Now()), f.EndWindow.Shift(time.Now()), f.NumPassengers, c)
				}(f, &wg, indx, provider)
			}

			if provider.IsAvailableRoute(f.SrcAirport, f.DestAirport) {
				wg.Add(1)
				go func(f config.FlightSearch, wg *sync.WaitGroup, indx int, provider providers.FlightProvider) {
					defer wg.Done()
					provider.GetOneWayFlights(f.DestAirport, f.SrcAirport, f.StartWindow.Shift(time.Now()), f.EndWindow.Shift(time.Now()), f.NumPassengers, c)
				}(f, &wg, indx, provider)
			}
		}
		//Wait for all one way flight requests to finish then close the channel
		go func(wg *sync.WaitGroup) {
			wg.Wait()
			close(c)
		}(&wg)
	}

	var wg sync.WaitGroup
	for _, ch := range chans {
		wg.Add(1)
		go func(c chan []*models.OneWayFlight) {
			defer wg.Done()
			for flights := range c {
				f.aggregate <- flights
			}
		}(ch)
	}

	wg.Wait()
	close(f.aggregate)
}
