package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"gitlab.com/buzz/ride/reqres"
)

var (
	requestRidePayload *reqres.RideRequestResponse
)

func TestRequestRideHTTPEndpoint(t *testing.T) {
	router := mux.NewRouter()

	router.Handle("/users/{userID}/request-ride", handleRideRequest(rideService{}))

	server := httptest.NewServer(router)

	requestURL := fmt.Sprintf("%s/users/%s/request-ride", server.URL, studentUserID)

	requestJSON := `{"latitude": "` + latitude + `", "longitude": "` + longitude + `"}`

	req, _ := http.NewRequest("POST", requestURL, strings.NewReader(requestJSON))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 201 {
		t.Errorf("Expected a 200 response status code but got: %d", resp.StatusCode)
	}
}

func TestGetRequestedRidesHTTPEndpoint(t *testing.T) {
	server := httptest.NewServer(handleGetAllRideRequest(rideService{}))

	getAllRequestURL := fmt.Sprintf("%s/requests", server.URL)

	resp, err := http.Get(getAllRequestURL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected a 200 response status code but got: %d", resp.StatusCode)
	}
}
