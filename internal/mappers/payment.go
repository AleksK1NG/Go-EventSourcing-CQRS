package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
	"github.com/AleksK1NG/es-microservice/internal/order/models"
	orderService "github.com/AleksK1NG/es-microservice/proto/order"
)

func PaymentFromProto(protoPayment *orderService.Payment) dto.Payment {
	return dto.Payment{
		PaymentID: protoPayment.GetID(),
		Timestamp: protoPayment.GetTimestamp().AsTime(),
	}
}

func PaymentResponseFromModel(payment models.Payment) dto.Payment {
	return dto.Payment{
		PaymentID: payment.PaymentID,
		Timestamp: payment.Timestamp,
	}
}
