package aggregate

import "github.com/pkg/errors"

var (
	ErrOrderAlreadyCompleted          = errors.New("Order already completed")
	ErrOrderAlreadyCanceled           = errors.New("Order is already canceled")
	ErrOrderMustBePaidBeforeDelivered = errors.New("Order must be paid before been delivered")
	ErrCancelReasonRequired           = errors.New("Cancel reason must be provided")
	ErrOrderAlreadyCancelled          = errors.New("order already cancelled")
	ErrAlreadyPaid                    = errors.New("already paid")
	ErrAlreadySubmitted               = errors.New("already submitted")
	ErrOrderNotPaid                   = errors.New("order not paid")
	ErrOrderNotFound                  = errors.New("order not found")
	ErrAlreadyCreated                 = errors.New("order with given id already created")
	ErrOrderShopItemsIsRequired       = errors.New("order shop items is required")
	ErrInvalidDeliveryAddress         = errors.New("Invalid delivery address")
)
