package grpc

import (
	"google.golang.org/grpc"
	"log"
	"sync"
)


var client FlightScraperClient

var conn *grpc.ClientConn

var wg sync.WaitGroup

func GetClient() FlightScraperClient {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	if client == nil {
		addr := "localhost:9999"
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		client = NewFlightScraperClient(conn)
	}

	return client
}

func CloseConnection() {
	_ = conn.Close()
}