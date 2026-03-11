package models

import "time"

type Booking struct {
	ID            string    `bson:"_id,omitempty" json:"id"`
	UserID        string    `bson:"user_id" json:"user_id"`
	Type          string    `bson:"type" json:"type"`
	Provider      string    `bson:"provider" json:"provider"`
	BookingRef    string    `bson:"booking_ref" json:"booking_ref"`
	PassengerName string    `bson:"passenger_name" json:"passenger_name"`
	FlightNumber  string    `bson:"flight_number" json:"flight_number"`
	Departure     string    `bson:"departure" json:"departure"`
	Arrival       string    `bson:"arrival" json:"arrival"`
	StartDate     time.Time `bson:"start_date" json:"start_date"`
	Status        string    `bson:"status" json:"status"`
	Confidence    float64   `bson:"confidence" json:"confidence"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
}