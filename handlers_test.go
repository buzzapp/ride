package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"gitlab.com/buzz/ride/reqres"
)

var (
	requestRideResponsePayload *reqres.RideRequestResponse
)

func TestGetAllRideRequestsHTTPEndpoint(t *testing.T) {
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

func TestRequestRideHTTPEndpoint(t *testing.T) {
	router := mux.NewRouter()

	router.Handle("/users/{userID}/request-ride", handleRideRequest(rideService{}))

	server := httptest.NewServer(router)

	requestURL := fmt.Sprintf("%s/users/%s/request-ride", server.URL, studentUserID)

	requestJSON := &reqres.RideRequestRequest{
		Latitude:  latitude,
		Longitude: longitude,
	}

	js, _ := json.Marshal(requestJSON)

	requestJSONReader := bytes.NewReader(js)

	req, _ := http.NewRequest("POST", requestURL, requestJSONReader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 201 {
		t.Errorf("Expected a 201 response status code but got: %d", resp.StatusCode)
	}

	getPayload(resp.Body)
}

func TestAcceptRideRequestHTTPEndpoint(t *testing.T) {
	// make sure the request id has been set
	if requestRideResponsePayload.Request.ID != "" {
		router := mux.NewRouter()

		router.Handle("/requests/{requestID}/accept", handleAcceptRideRequest(rideService{}))

		server := httptest.NewServer(router)

		acceptRequestURL := fmt.Sprintf("%s/requests/%s/accept", server.URL, requestRideResponsePayload.Request.ID)

		req, _ := http.NewRequest("POST", acceptRequestURL, nil)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("Expected a 200 response status code but got: %d", resp.StatusCode)
		}
	}
}

func getPayload(respBody io.Reader) {
	var payload = &reqres.RideRequestResponse{}
	if err := json.NewDecoder(respBody).Decode(&payload); err != nil {
		log.Fatal("error decoding request response payload ", err)
	}

	requestRideResponsePayload = payload
}

func tearDown() {
	//Grab a copy of our session
	session, err := getSession()
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	//Get our collection of applications
	db := session.DB("buzz-test-ride")
	collection := db.C("requests")

	//remove our applications from the collection
	_, err = collection.RemoveAll(bson.M{"latitude": requestRideResponsePayload.Request.Latitude})
	if err != nil {
		fmt.Println("error removing test request ", err)
	}
}
