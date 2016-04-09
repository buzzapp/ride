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
	Version = "0.1.0"
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

	httpAddress := ":5000"

	go func() {
		l.Info("Establishing HTTP Bindings", "Main", "addr", httpAddress, "transport", "HTTP/JSON")

		// Create a new mux router
		router := mux.NewRouter()

		const GetAllRideRequestURL = "/requests"
		router.Handle(GetAllRideRequestURL, handleGetAllRideRequest(service)).Methods("GET")
		l.Info("New Handler", "Main", "path", GetAllRideRequestURL, "type", "GET")

		// register our router and start the server
		handler := cors.Default().Handler(router)
		http.Handle("/", router)
		errc <- http.ListenAndServe(httpAddress, handler)
	}()

	fmt.Println("Fatal Error", "Main", <-errc)
}