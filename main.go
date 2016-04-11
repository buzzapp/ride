package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/forestgiant/semver"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	//Version represents the semantic version of this service/api
	Version     = "0.1.0"
	httpAddress = ":8001"
)

func main() {
	// Setup Semantic Version flags
	err := semver.SetVersion(Version)
	if err != nil {
		log.Fatal(err)
	}

	// Check for command line configuration flags
	var (
		logPathUsage = "Path to the service logs."
		logPathPtr   = flag.String("logpath", "", logPathUsage)
	)
	flag.Parse()

	if len(*logPathPtr) == 0 {
		log.Fatal("You must provide a path where log files can be stored.")
	}

	l := getLogger(*logPathPtr)

	// `package log` domain
	l.Info("Initializing app.", "Main")

	// Mechanical stuff
	errc := make(chan error)
	go func() {
		errc <- interrupt()
	}()

	// Define our app service
	var service RideService
	service = rideService{}

	go func() {
		l.Info("Establishing HTTP Bindings", "Main", "addr", httpAddress, "transport", "HTTP/JSON")

		// Create a new mux router
		router := mux.NewRouter()

		const AcceptRideRequestPath = "/requests/{requestID}/accept"
		router.Handle(AcceptRideRequestPath, handleAcceptRideRequest(service)).Methods("POST")
		l.Info("New Handler", "Main", "path", AcceptRideRequestPath, "type", "POST")

		const GetRequestedRidesPath = "/requests"
		router.Handle(GetRequestedRidesPath, handleGetAllRideRequest(service)).Methods("GET")
		l.Info("New Handler", "Main", "path", GetRequestedRidesPath, "type", "GET")

		const RideRequestPath = "/users/{userID}/requests"
		router.Handle(RideRequestPath, handleRideRequest(service)).Methods("POST")
		l.Info("New Handler", "Main", "path", RideRequestPath, "type", "POST")

		// register our router and start the server
		http.Handle("/", router)
		c := cors.New(cors.Options{
			AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "Authorization"},
		})
		handler := c.Handler(router)
		errc <- http.ListenAndServe(httpAddress, handler)
	}()

	fmt.Println("Fatal Error", "Main", <-errc)
}
