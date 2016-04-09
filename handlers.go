package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"gitlab.com/buzz/ride/model"
	"gitlab.com/buzz/ride/reqres"
)

func handleAcceptRideRequest(svc RideService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := mux.Vars(r)["requestID"]

		// do validation
		if requestID == "" {
			respondWithError("Invalid request", errors.New("no request id passed in"), w, http.StatusBadGateway)
			return
		}

		if err := svc.AcceptRideRequest(requestID); err != nil {
			respondWithError("Unable to request ride", err, w, http.StatusInternalServerError)
			return
		}

		// Generate our response
		resp := reqres.MessageResponse{Message: "ride " + requestID + " was accepted"}

		// Marshal up the json response
		js, err := json.Marshal(resp)
		if err != nil {
			respondWithError("unable to marshal json response", err, w, http.StatusInternalServerError)
			return
		}

		// Return the response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})
}

func handleRideRequest(svc RideService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["userID"]

		if userID == "" {
			respondWithError("Invalid request", errors.New("no user id passed in"), w, http.StatusBadGateway)
			return
		}

		var payload = &reqres.RideRequestRequest{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			respondWithError("unable to decode json request", err, w, http.StatusInternalServerError)
			return
		}

		// do validation

		request, err := svc.RequestRide(userID, payload.Latitude, payload.Longitude)
		if err != nil {
			respondWithError("Unable to request ride", err, w, http.StatusInternalServerError)
			return
		}

		// Generate our response
		resp := reqres.RideRequestResponse{Request: request}

		// Marshal up the json response
		js, err := json.Marshal(resp)
		if err != nil {
			respondWithError("unable to marshal json response", err, w, http.StatusInternalServerError)
			return
		}

		// Return the response
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})
}

func handleGetAllRideRequest(svc RideService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.RequestURI)
		if err != nil {
			respondWithError("Unable to parse url", err, w, http.StatusInternalServerError)
			return
		}

		acceptedStatus := u.Query().Get("accepted")

		var rides []model.Request

		if acceptedStatus == "" {
			filter := make(map[string]interface{})
			rides, err = svc.GetAllRideRequests(filter)
			if err != nil {
				respondWithError("Unable to get requested rides", err, w, http.StatusInternalServerError)
				return
			}
		} else {
			filter := make(map[string]interface{})
			if acceptedStatus == "true" {
				filter["accepted"] = true
			} else {
				filter["accepted"] = false
			}
			rides, err = svc.GetAllRideRequests(filter)
			if err != nil {
				respondWithError("Unable to get requested rides", err, w, http.StatusInternalServerError)
				return
			}
		}

		// Generate our response
		resp := reqres.GetAllRideRequestResponse{Requests: rides}

		// Marshal up the json response
		js, err := json.Marshal(resp)
		if err != nil {
			respondWithError("unable to marshal json response", err, w, http.StatusInternalServerError)
			return
		}

		// Return the response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})
}

// Helper function to return a json error message
func respondWithError(msg string, err error, w http.ResponseWriter, status int) {
	errMsg := reqres.ErrorResponse{Message: msg + ": " + err.Error()}

	js, err := json.Marshal(errMsg)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}
