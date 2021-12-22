package southwest

import (
	"context"
	"fmt"
	pb "github.com/gabriel-flynn/Cheap-Flight-Finder/server/grpc"
	"log"
	"net/http"
	"strconv"
)

func refreshAuthHeaders()  {
	fmt.Println("REFRESHING AUTH HEADERS!!!")
	client := pb.GetClient()

	req := pb.Empty{
	}
	resp, err := client.GetSouthwestHeaders(context.Background(), &req)

	//fmt.Println(resp.Headers)
	if err != nil {
		log.Fatal(err)
	}
	ApiInfo.Lock()
	ApiInfo.Headers = resp.Headers
	ApiInfo.Unlock()
}


func oneWayFlightRequest(srcAirport, destAirport, departureDate string, numPassengers int, roundTripFlight bool) *http.Request {

	req, _ := http.NewRequest("GET", MobileFlightUrl, nil)
	ApiInfo.RLock()
	headers := ApiInfo.Headers

	req.Header.Set("EE30zvQLWf-a", headers["EE30zvQLWf-a"])
	req.Header.Set("EE30zvQLWf-b", headers["EE30zvQLWf-b"])
	req.Header.Set("EE30zvQLWf-c", headers["EE30zvQLWf-c"])
	req.Header.Set("EE30zvQLWf-d", headers["EE30zvQLWf-d"])
	req.Header.Set("EE30zvQLWf-f", headers["EE30zvQLWf-f"])
	req.Header.Set("EE30zvQLWf-z", headers["EE30zvQLWf-z"])
	req.Header.Set("x-app-version", headers["x-app-version"])
	req.Header.Set("X-Requested-With", headers["X-Requested-With"])
	req.Header.Set("Accept", headers["Accept"])
	req.Header.Set("X-Channel-ID", headers["X-Channel-ID"])
	req.Header.Set("x-app-version", headers["x-app-version"])
	req.Header.Set("X-API-Key", headers["X-API-Key"])
	req.Header.Set("X-User-Experience-ID", headers["X-User-Experience-ID"])
	req.Header.Set("User-Agent", headers["User-Agent"])

	//for k, v := range config2.ApiInfo.Headers {
	//	fmt.Printf("%s: %s\n", k, v)
	//	req.Header.Set(k, v)
	//}


	defer ApiInfo.RUnlock()

	q := req.URL.Query()
	q.Add("origination-airport", srcAirport)
	q.Add("destination-airport", destAirport)
	q.Add("departure-date", departureDate)
	q.Add("number-adult-passengers", strconv.FormatInt(int64(numPassengers), 10))
	q.Add("currency", "USD")

	fmt.Println(q.Get("number-adult-passengers"))
	req.URL.RawQuery = q.Encode()

	return req
}
