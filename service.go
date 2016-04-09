package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pusher/pusher-http-go"

	"gitlab.com/buzz/ride/model"
	"gitlab.com/buzz/ride/reqres"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// SecretKey is the key to hash the JWT token
	SecretKey = "33266AB738F764C2A3DD5D8F38336"
)

// RideService needs a description
type RideService interface {
	AcceptRideRequest(requestID string) error
	GetAllRideRequests(filters map[string]interface{}) ([]model.Request, error)
	GetRideRequestByID(id string) (*model.Request, error)
	RequestRide(userID string, latitude, longitude float32) (*model.Request, error)
}

type rideService struct{}

func (svc rideService) AcceptRideRequest(requestID string) error {
	acceptedRequest, err := svc.GetRideRequestByID(requestID)
	if err != nil {
		return errors.New("Error getting ride " + err.Error())
	}

	acceptedRequest.Accepted = true

	//Grab a copy of our session
	session, err := getSession()
	if err != nil {
		return err
	}
	defer session.Close()

	//Get our collection of requests
	db := session.DB("buzz-test-ride")
	collection := db.C("requests")

	//update our request
	err = collection.Update(bson.M{"_id": requestID}, acceptedRequest)
	if err != nil {
		return err
	}

	/************************* CLEAN UP *****************************/
	requests, err := svc.GetAllRideRequests(bson.M{"accepted": false})
	if err != nil {
		fmt.Println(err)
	} else {
		// send pusher notification
		client := pusher.Client{
			AppId:  "196361",
			Key:    "524d992dcbb206bbae9a",
			Secret: "e43cf33b3a27d6298ce7",
			Secure: true,
		}

		_, err = client.Trigger("request_channel", "ride-accepted", requests)
		if err != nil {
			return err
		}
	}
	/****************************************************************/

	return nil
}

func (svc rideService) GetAllRideRequests(filters map[string]interface{}) ([]model.Request, error) {
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

func (rideService) GetRideRequestByID(id string) (*model.Request, error) {
	//Grab a copy of our session
	session, err := getSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	//Get our collection of applications
	db := session.DB("buzz-test-ride")
	collection := db.C("requests")

	//Get our requests from the collection
	var retrievedRequest *model.Request
	err = collection.Find(bson.M{"_id": id}).Sort("-timestamp").One(&retrievedRequest)
	if err != nil {
		return nil, err
	}

	return retrievedRequest, nil
}

func (rideService) RequestRide(userID string, latitude, longitude float32) (*model.Request, error) {
	// Get a user
	userURL := fmt.Sprintf("http://localhost:8000/users/%s", userID)
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

	/************************* CLEAN UP *****************************/

	// send pusher notification
	client := pusher.Client{
		AppId:  "196361",
		Key:    "524d992dcbb206bbae9a",
		Secret: "e43cf33b3a27d6298ce7",
		Secure: true,
	}

	_, err = client.Trigger("request_channel", "ride-requested", request)
	if err != nil {
		return nil, err
	}

	/****************************************************************/

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
