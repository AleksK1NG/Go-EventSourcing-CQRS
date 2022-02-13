package dto

type ChangeDeliveryAddressReqDto struct {
	DeliveryAddress string `json:"deliveryAddress" bson:"deliveryAddress,omitempty" validate:"required"`
}
