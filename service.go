package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pusher/pusher-http-go"

	"gitlab.com/buzz/ride/model"
	"gitlab.com/buzz/ride/reqres"
	stela "gitlab.fg/go/stela/api"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// RideService needs a description
type RideService interface {
	GetAllRideRequests(filters map[string]interface{}) ([]model.Request, error)
	RequestRide(userID, latitude, longitude string) (*model.Request, error)
	UpdateRideRequest(updatedRequest *model.UpdateRequest) (*model.Request, error)
}

type rideService struct{}

func (rideService) GetAllRideRequests(filters map[string]interface{}) ([]model.Request, error) {
	//Grab a copy of our session
	session, err := getSession()
	if err != nil {
		return []model.Request{}, err
	}
	defer session.Close()

	//Get our collection of applications
	db := session.DB("buzz-test-ride")
	collection := db.C("requests")

	//Get our requests from the collection
	var retrievedRequests []model.Request
	err = collection.Find(filters).Sort("-timestamp").All(&retrievedRequests)
	if err != nil {
		return []model.Request{}, err
	}

	return retrievedRequests, nil
}

func (rideService) RequestRide(userID, latitude, longitude string) (*model.Request, error) {
	// Create stela client
	sclient, err := stela.NewClient("localhost:9000")
	if err != nil {
		return nil, err
	}

	// Discover all endpoints registered with that name
	// Typically done from another service
	service, err := sclient.DiscoverOne("user.service.buzz")
	if err != nil {
		return nil, err
	}

	// Get a user
	userURL := fmt.Sprintf("http://%s:%d/users/%s", service.Address, service.Port, userID)
	resp, err := http.Get(userURL)
	if err != nil {
		return nil, err
	}

	var userPayload = &reqres.GetUserByIDResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&userPayload); err != nil {
		return nil, err
	}

	request := &model.Request{
		ID:        bson.NewObjectId().Hex(),
		User:      userPayload.User,
		Latitude:  latitude,
		Longitude: longitude,
		Accepted:  false,
		CreatedAt: time.Now().Unix(),
	}

	//Grab a copy of our session
	session, err := getSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	//Get our collection of applications
	db := session.DB("buzz-test-ride")
	collection := db.C("requests")

	if err := collection.Insert(request); err != nil {
		return nil, err
	}

	// send pusher notification
	client := pusher.Client{
		AppId:  "196361",
		Key:    "524d992dcbb206bbae9a",
		Secret: "e43cf33b3a27d6298ce7",
		Secure: true,
	}

	client.Trigger("test_channel", "my_event", request)

	return request, nil
}

func (rideService) UpdateRideRequest(updatedRequest *model.UpdateRequest) (*model.Request, error) {
	updatedReq := &model.Request{
		ID:        updatedRequest.ID,
		Latitude:  updatedRequest.Latitude,
		Longitude: updatedRequest.Longitude,
		Accepted:  updatedRequest.Accepted,
		CreatedAt: updatedRequest.CreatedAt,
		User:      updatedRequest.User,
		UpdatedAt: updatedRequest.UpdatedAt,
	}

	//Grab a copy of our session
	session, err := getSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	//Get our collection of applications
	db := session.DB("buzz-test-ride")
	collection := db.C("requests")

	//Insert our application
	err = collection.Update(bson.M{"_id": updatedRequest.ID}, updatedReq)
	if err != nil {
		return nil, err
	}

	return updatedReq, nil
}

var globalSession *mgo.Session

func getSession() (*mgo.Session, error) {
	//Establish our database connection
	if globalSession == nil {
		var err error
		globalSession, err = mgo.Dial(":27017")
		if err != nil {
			return nil, err
		}

		//Optional. Switch the session to a monotonic behavior.
		globalSession.SetMode(mgo.Monotonic, true)
	}

	return globalSession.Copy(), nil
}
