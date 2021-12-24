package frontier

import (
	"net/http"
)

func getRoutesRequest() *http.Request {
	req, _ := http.NewRequest("GET", FrontierRoutesUrl, nil)
	return req
}
