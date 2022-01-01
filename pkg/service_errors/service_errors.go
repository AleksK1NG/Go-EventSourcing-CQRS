package serviceErrors

import (
	"github.com/pkg/errors"
)

const (
	ErrMsgMongoCollectionAlreadyExists = "Collection already exists"
)

var (
	ErrAlreadyCreatedOrCancelled = errors.New("order created or cancelled")
	ErrAlreadyPaid               = errors.New("already paid")
	ErrAlreadySubmitted          = errors.New("already submitted")
	ErrAlreadyCreated            = errors.New("already created")
	ErrOrderNotPaid              = errors.New("order not paid")
	ErrSubscriptionDropped       = errors.New("Subscription Dropped")
	ErrOrderNotFound             = errors.New("order not found")
)
