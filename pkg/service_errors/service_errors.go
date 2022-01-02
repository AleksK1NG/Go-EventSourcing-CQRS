package serviceErrors

import (
	"github.com/pkg/errors"
)

const (
	ErrMsgMongoCollectionAlreadyExists = "Collection already exists"
	ErrMsgAlreadyExists                = "already exists"
)

var (
	ErrAlreadyCreatedOrCancelled = errors.New("order created or cancelled")
	ErrAlreadyPaid               = errors.New("already paid")
	ErrAlreadySubmitted          = errors.New("already submitted")
	ErrOrderNotPaid              = errors.New("order not paid")
	ErrOrderNotFound             = errors.New("order not found")
)
