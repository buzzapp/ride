package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/buzzapp/ride/model"

	"gopkg.in/mgo.v2/bson"
)

const (
	studentUserID = "5707fa8ae5b07e2f4f91e883"
)

var (
	fromAddressTest = model.Address{Latitude: 38.209265, Longitude: -85.763981}
	toAddressTest   = model.Address{Latitude: 38.2543950, Longitude: -85.6875660}
)

var (
	rideSVC            RideService
	serviceRequestRide *model.Request
)

func TestMain(m *testing.M) {
	rideSVC = rideService{}

	result := m.Run()

	tearDown()

	os.Exit(result)
}

func TestRequestRide(t *testing.T) {
	request, err := rideSVC.RequestRide(studentUserID, fromAddressTest, toAddressTest)
	if err != nil {
		t.Error(err)
	}

	if request.User.ID != studentUserID {
		t.Error("Request student id does not match")
	}

	if request.FromAddress != fromAddressTest {
		t.Error("Request latitude does not match")
	}

	if request.ToAddress != toAddressTest {
		t.Error("Request longitude does not match")
	}

	serviceRequestRide = request
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

func TestAcceptRideRequest(t *testing.T) {
	if err := rideSVC.AcceptRideRequest(serviceRequestRide.ID); err != nil {
		t.Error(err)
	}
}

func TestGetRideRequestByID(t *testing.T) {
	ride, err := rideSVC.GetRideRequestByID(serviceRequestRide.ID)
	if err != nil {
		t.Error(err)
	}

	if ride.ID != serviceRequestRide.ID {
		t.Errorf("error getting app by id expecting app with id of %s but got %s", serviceRequestRide.ID, ride.ID)
	}
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
	_, err = collection.RemoveAll(bson.M{"from_address": serviceRequestRide.FromAddress})
	if err != nil {
		fmt.Println("error removing test request ", err)
	}
}
