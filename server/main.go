package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/config"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/coordinator"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/datastore"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/models"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	flightConfig := config.LoadConfiguration("config.json")
	c := make(chan []*models.OneWayFlight)

	t := time.Now()
	//Start querying endpoints for flights
	go coordinator.NewFlightProviderCoordinator(c).Start()

	var flights []*models.OneWayFlight

	file, err := os.Create("flights.csv")
	if err != nil {
		log.Fatal("Couldn't create file.", err)
	}

	writer := csv.NewWriter(file)
	writer.Write(models.GetOneWayFlightHeading())

	var wg sync.WaitGroup
	for msg := range c {
		flights = append(flights, msg...)
		if len(flights) > 500 {
			err := writer.WriteAll(models.BatchGetOneWayFlightRecord(flights))
			if err != nil {
				log.Fatal("Error while writing csv:", err)
			}

			if flightConfig.SaveToDynamo {
				//Save to dynamodb
				wg.Add(1)
				flightsCopy := make([]*models.OneWayFlight, len(flights))
				copy(flightsCopy, flights)
				go func() {
					datastore.SaveOneWayFlightBatch(flightsCopy)
					wg.Done()
				}()
			}
			flights = []*models.OneWayFlight{}
		}
	}

	if len(flights) != 0 {
		err := writer.WriteAll(models.BatchGetOneWayFlightRecord(flights))

		if flightConfig.SaveToDynamo {
			//Save to dynamo
			wg.Add(1)
			copySlice := make([]*models.OneWayFlight, len(flights))
			copy(copySlice, flights)
			go func() {
				datastore.SaveOneWayFlightBatch(copySlice)
				wg.Done()
			}()
			if err != nil {
				log.Fatal("Error while writing csv:", err)
			}
		}
		flights = []*models.OneWayFlight{}
	}
	file.Close()

	fmt.Println(time.Since(t))
	if flightConfig.SaveToDynamo {
		fmt.Print("Waiting for flights to save to DynamoDB")
	}
	wg.Wait()
}