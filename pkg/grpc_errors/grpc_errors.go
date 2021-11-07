package grpcErrors

import (
	"context"
	"database/sql"
	"github.com/AleksK1NG/es-microservice/pkg/constants"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

var (
	ErrNoCtxMetaData = errors.New("No ctx metadata")
)

//ErrResponse get gRPC error response
func ErrResponse(err error) error {
	return status.Error(GetErrStatusCode(err), err.Error())
}

//GetErrStatusCode get error status code from error
func GetErrStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return codes.NotFound
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	case errors.Is(err, ErrNoCtxMetaData):
		return codes.Unauthenticated
	case CheckErrMessage(err, constants.Validate):
		return codes.InvalidArgument
	case CheckErrMessage(err, constants.Redis):
		return codes.NotFound
	case CheckErrMessage(err, constants.FieldValidation):
		return codes.InvalidArgument
	case CheckErrMessage(err, constants.RequiredHeaders):
		return codes.Unauthenticated
	case CheckErrMessage(err, constants.Base64):
		return codes.InvalidArgument
	case CheckErrMessage(err, constants.Unmarshal):
		return codes.InvalidArgument
	case CheckErrMessage(err, constants.Uuid):
		return codes.InvalidArgument
	case CheckErrMessage(err, constants.Cookie):
		return codes.Unauthenticated
	case CheckErrMessage(err, constants.Token):
		return codes.Unauthenticated
	case CheckErrMessage(err, constants.Bcrypt):
		return codes.InvalidArgument
	}
	return codes.Internal
}

func CheckErrMessage(err error, msg string) bool {
	return strings.Contains(strings.ToLower(err.Error()), msg)
}
