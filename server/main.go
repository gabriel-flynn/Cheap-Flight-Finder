package main

import (
	"context"
	"fmt"
	pb "github.com/gabriel-flynn/Cheap-Flight-Finder/server/grpc"
	"google.golang.org/grpc"
	"log"
)

func main() {
	// config := config.LoadConfiguration("config.json")
	// c := make(chan []*models.OneWayFlight)

	// addr := "localhost:9999"
	// conn, err := grpc.Dial(addr, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer conn.Close()
	// //c := pb.NewFlightScraperClient(conn)


	// 		if config.SaveToDynamo {
	// 			//Save to dynamodb
	// 			wg.Add(1)
	// 			flightsCopy := make([]*models.OneWayFlight, len(flights))
	// 			copy(flightsCopy, flights)
	// 			go func() {
	// 				database.SaveOneWayFlightBatch(flightsCopy)
	// 				wg.Done()
	// 			}()
	// 		}
	// 		flights = []*models.OneWayFlight{}
	// 	}
	// }

	// if len(flights) != 0 {
	// 	err := writer.WriteAll(models.BatchGetOneWayFlightRecord(flights))

	// 	if config.SaveToDynamo {
	// 		//Save to dynamo
	// 		wg.Add(1)
	// 		copySlice := make([]*models.OneWayFlight, len(flights))
	// 		copy(copySlice, flights)
	// 		go func() {
	// 			database.SaveOneWayFlightBatch(copySlice)
	// 			wg.Done()
	// 		}()
	// 		if err != nil {
	// 			log.Fatal("Error while writing csv:", err)
	// 		}
	// 	}
	// 	flights = []*models.OneWayFlight{}
	// }
	// file.Close()
	// if config.SaveToDynamo {
	// 	fmt.Print("Waiting for flights to save to DynamoDB")
	// }
	// wg.Wait()

	addr := "localhost:9999"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	//c := pb.NewFlightScraperClient(conn)
client := pb.NewFlightScraperClient(conn)
	req := pb.Empty{
	}

	resp, err := client.GetSouthwestHeaders(context.Background(), &req)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range resp.Headers {
		fmt.Println(k, "value is", v)
	}
}
