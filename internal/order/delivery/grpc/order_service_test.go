package grpc

import (
	"context"
	"github.com/AleksK1NG/es-microservice/pkg/interceptors"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	orderService "github.com/AleksK1NG/es-microservice/proto/order"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"sync"
	"testing"
	"time"
)

func NewOrderServiceConn(ctx context.Context, im interceptors.InterceptorManager) (*grpc.ClientConn, error) {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(500 * time.Millisecond)),
		grpc_retry.WithCodes(codes.NotFound, codes.Aborted),
		grpc_retry.WithMax(5),
	}

	orderGrpcConn, err := grpc.DialContext(
		ctx,
		":5001",
		grpc.WithUnaryInterceptor(im.ClientRequestLoggerInterceptor()),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
	)
	if err != nil {
		return nil, errors.Wrap(err, "grpc.DialContext")
	}

	return orderGrpcConn, nil
}

func TestOrderGrpcService_UpdateOrder(t *testing.T) {
	appLogger := logger.NewAppLogger(&logger.Config{LogLevel: "debug", DevMode: false, Encoder: "json"})
	appLogger.InitLogger()
	appLogger.WithName("OrderService")

	im := interceptors.NewInterceptorManager(appLogger)

	orderServiceConn, err := NewOrderServiceConn(context.Background(), im)
	if err != nil {
		return
	}
	defer orderServiceConn.Close()

	client := orderService.NewOrderServiceClient(orderServiceConn)
	aggregateID := "ead2b45d-cf1c-4500-9cae-a556bf3e4687"

	wg := &sync.WaitGroup{}
	for i := 0; i <= 60; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			result, err := client.UpdateOrder(context.Background(), &orderService.UpdateOrderReq{
				AggregateID: aggregateID,
				ItemID:      []string{uuid.NewV4().String(), uuid.NewV4().String(), uuid.NewV4().String()},
			})
			if err != nil {
				appLogger.WarnMsg("client.UpdateOrder", err)
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			appLogger.Infof("result: %s", result.String())

			res, err := client.PayOrder(context.Background(), &orderService.PayOrderReq{AggregateID: aggregateID})
			if err != nil {
				appLogger.WarnMsg("client.PayOrder", err)
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			appLogger.Infof("res: %s", res.String())

			for i := 0; i < 5; i++ {
				result, err = client.UpdateOrder(context.Background(), &orderService.UpdateOrderReq{
					AggregateID: aggregateID,
					ItemID:      []string{uuid.NewV4().String(), uuid.NewV4().String(), uuid.NewV4().String()},
				})
				if err != nil {
					appLogger.WarnMsg("client.UpdateOrder", err)
				}
				require.NoError(t, err)
				require.NotNil(t, result)
				appLogger.Infof("result: %s", result.String())
			}

			resp, err := client.SubmitOrder(context.Background(), &orderService.SubmitOrderReq{AggregateID: aggregateID})
			if err != nil {
				appLogger.WarnMsg("client.SubmitOrder", err)
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			appLogger.Infof("resp: %s", resp.String())

			result, err = client.UpdateOrder(context.Background(), &orderService.UpdateOrderReq{
				AggregateID: aggregateID,
				ItemID:      []string{uuid.NewV4().String(), uuid.NewV4().String(), uuid.NewV4().String()},
			})
			if err != nil {
				appLogger.WarnMsg("client.UpdateOrder", err)
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			appLogger.Infof("result: %s", result.String())

		}(wg)
	}
	wg.Wait()
	//for i := 0; i <= 20; i++ {
	//	result, err := client.UpdateOrder(context.Background(), &orderService.UpdateOrderReq{
	//		AggregateID: "c729afba-2b5b-4cea-b9c3-3cd140be6164",
	//		ItemID:      []string{uuid.NewV4().String(), uuid.NewV4().String(), uuid.NewV4().String()},
	//	})
	//	if err != nil {
	//		appLogger.WarnMsg("client.UpdateOrder", err)
	//	}
	//	require.NoError(t, err)
	//	require.NotNil(t, result)
	//	appLogger.Infof("result: %s", result.String())
	//}

	appLogger.Infof("Success =D")
}
