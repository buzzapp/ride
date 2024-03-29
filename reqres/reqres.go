package reqres

import "github.com/buzzapp/ride/model"

// RideRequestRequest describes a request for a new ride request
type RideRequestRequest struct {
	FromAddress model.Address `json:"from_address"`
	ToAddress   model.Address `json:"to_address"`
}

// RideRequestResponse describes a response for a new ride request
type RideRequestResponse struct {
	Request *model.Request `json:"request"`
}

// GetUserByIDResponse describes a resposne for getting a user
type GetUserByIDResponse struct {
	User model.User `json:"user"`
}

// GetAllRideRequestResponse describes a resposne for getting all rides
type GetAllRideRequestResponse struct {
	Requests []model.Request `json:"requests"`
}

/*****************************/
/* GENERIC RESPONSES */
/*****************************/

// ErrorResponse describes a response for when there is an error
type ErrorResponse struct {
	Message string `json:"message"`
}

// MessageResponse describes a message JSON response
type MessageResponse struct {
	Message string `json:"message"`
}
