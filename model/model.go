package model

// Request represnts a ride request placed by a student user
type Request struct {
	ID        string `bson:"_id" json:"id"`
	Latitude  string `bson:"latitude" json:"latitude"`
	Longitude string `bson:"longitude" json:"longitude"`
	Accepted  bool   `bson:"accepted" json:"accepted"`
	Timestamp int64  `bson:"timestamp" json:"timestamp"`
	User      User   `bson:"user" json:"user"`
}

// User describes the properties of a user
type User struct {
	ID        string `bson:"_id" json:"id"`
	Email     string `bson:"email" json:"email"`
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	Password  string `bson:"password" json:"password"`
	Role      string `bson:"role" json:"role"`
	Username  string `bson:"username" json:"username"`
	Timestamp int64  `bson:"timestamp" json:"timestamp"`
}
