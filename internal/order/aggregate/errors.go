package aggregate

import "github.com/pkg/errors"

var (
	ErrOrderAlreadyDelivered          = errors.New("Order already delivered")
	ErrOrderAlreadyCanceled           = errors.New("Order is already canceled")
	ErrOrderMustBePaidBeforeDelivered = errors.New("Order must be paid before been delivered")
	ErrCancelReasonRequired           = errors.New("Cancel reason must be provided")
	ErrAlreadyCreatedOrCancelled      = errors.New("order created or cancelled")
	ErrAlreadyPaid                    = errors.New("already paid")
	ErrAlreadySubmitted               = errors.New("already submitted")
	ErrOrderNotPaid                   = errors.New("order not paid")
	ErrOrderNotFound                  = errors.New("order not found")
	ErrAlreadyCreated                 = errors.New("order with given id already created")
	ErrOrderShopItemsIsRequired       = errors.New("order shop items is required")
	ErrInvalidDeliveryAddress         = errors.New("Invalid delivery address")
)
