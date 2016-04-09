package main

import (
	"log"
	"os"
	"testing"
	"time"

	"gitlab.com/buzz/ride/model"

	"gopkg.in/mgo.v2/bson"
)

const (
	studentUserID = "5707fa8ae5b07e2f4f91e883"
	latitude      = "38.253238677978"
	longitude     = "-85.662582397461"
)

var (
	rideSVC   RideService
	requestID string
	accepted  = true
	id        string
	user      model.User
	createdAt int64
)

func TestMain(m *testing.M) {
	rideSVC = rideService{}

	result := m.Run()

	tearDown()

	os.Exit(result)
}

func TestRequestRide(t *testing.T) {
	request, err := rideSVC.RequestRide(studentUserID, latitude, longitude)
	if err != nil {
		t.Error(err)
	}

	if request.User.ID != studentUserID {
		t.Error("Request student id does not match")
	}

	if request.Latitude != latitude {
		t.Error("Request latitude does not match")
	}

	if request.Longitude != longitude {
		t.Error("Request longitude does not match")
	}

	requestID = request.ID
	user = request.User
	createdAt = request.CreatedAt
}

func TestGetAllRideRequests(t *testing.T) {
	filter := make(map[string]interface{})
	_, err := rideSVC.GetAllRideRequests(filter)
	if err != nil {
		t.Error(err)
	}

	// check filer works
	filter["accepted"] = true
	rr, err := rideSVC.GetAllRideRequests(filter)
	if err != nil {
		t.Error(err)
	}

	for _, r := range rr {
		if r.Accepted != true {
			t.Error("filter not working")
		}
	}

	filter["accepted"] = false
	rr, err = rideSVC.GetAllRideRequests(filter)
	if err != nil {
		t.Error(err)
	}

	for _, r := range rr {
		if r.Accepted != false {
			t.Error("filter not working")
		}
	}
}

func TestUpdateRequest(t *testing.T) {
	updatedReq := &model.UpdateRequest{
		ID:        requestID,
		Latitude:  latitude,
		Longitude: longitude,
		Accepted:  accepted,
		CreatedAt: createdAt,
		UpdatedAt: time.Now().Unix(),
	}

	ride, err := rideSVC.UpdateRideRequest(updatedReq)
	if err != nil {
		t.Error(err)
	}

	if ride.Accepted != accepted {
		t.Error("not accepted")
	}
}

func tearDown() {
	//Grab a copy of our session
	session, err := getSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	//Get our collection of applications
	db := session.DB("buzz-test-ride")
	collection := db.C("requests")

	//remove our applications from the collection
	err = collection.Remove(bson.M{"latitude": latitude})
	if err != nil {
		log.Fatal(err)
	}
}
