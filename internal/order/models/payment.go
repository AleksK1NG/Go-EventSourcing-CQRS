package models

import (
	"fmt"
	"time"
)

type Payment struct {
	PaymentID string    `json:"paymentID" bson:"paymentID,omitempty" validate:"required"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp,omitempty" validate:"required"`
}

func (p *Payment) String() string {
	return fmt.Sprintf("PaymentID: {%s}, Timestamp: {%s}", p.PaymentID, p.Timestamp.UTC().String())
}
