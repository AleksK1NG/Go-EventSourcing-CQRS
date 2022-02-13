package dto

import "time"

type Payment struct {
	PaymentID string    `json:"paymentID" bson:"paymentID,omitempty" validate:"required"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp,omitempty" validate:"required"`
}
