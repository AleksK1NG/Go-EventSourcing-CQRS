package dto

type CancelOrderReqDto struct {
	CancelReason string `json:"cancelReason" validate:"required"`
}
