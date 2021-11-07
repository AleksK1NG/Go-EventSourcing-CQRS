package serviceErrors

import (
	"github.com/pkg/errors"
)

var (
	ErrAlreadyCreatedOrCancelled = errors.New("order created or cancelled")
	ErrAlreadyPaid               = errors.New("already paid")
	ErrAlreadySubmitted          = errors.New("already submitted")
	ErrAlreadyCreated            = errors.New("already created")
	ErrOrderNotPaid              = errors.New("order not paid")
	ErrSubscriptionDropped       = errors.New("Subscription Dropped")
)
