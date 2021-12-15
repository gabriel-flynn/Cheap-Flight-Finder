package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/config"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/coordinator"
	database "github.com/gabriel-flynn/Cheap-Flight-Finder/server/datastore"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/models"
	"log"
	"os"
	"sync"
)

func main() {
	config := config.LoadConfiguration("config.json")
	c := make(chan []*models.OneWayFlight)

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

			if config.SaveToDynamo {
				//Save to dynamodb
				wg.Add(1)
				flightsCopy := make([]*models.OneWayFlight, len(flights))
				copy(flightsCopy, flights)
				go func() {
					database.SaveOneWayFlightBatch(flightsCopy)
					wg.Done()
				}()
			}
			flights = []*models.OneWayFlight{}
		}
	}

	if len(flights) != 0 {
		err := writer.WriteAll(models.BatchGetOneWayFlightRecord(flights))

		if config.SaveToDynamo {
			//Save to dynamo
			wg.Add(1)
			copySlice := make([]*models.OneWayFlight, len(flights))
			copy(copySlice, flights)
			go func() {
				database.SaveOneWayFlightBatch(copySlice)
				wg.Done()
			}()
			if err != nil {
				log.Fatal("Error while writing csv:", err)
			}
		}
		flights = []*models.OneWayFlight{}
	}
	file.Close()
	if config.SaveToDynamo {
		fmt.Print("Waiting for flights to save to DynamoDB")
	}
	wg.Wait()
}
