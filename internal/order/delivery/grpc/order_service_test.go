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
	aggregateID := "22699160-548f-4422-9d70-6484e2a612a8"

	res, err := client.CreateOrder(context.Background(), &orderService.CreateOrderReq{
		AccountEmail: "alexander.bryksin@yandex.ru",
		ShopItems: []*orderService.ShopItem{
			{
				ID:          uuid.NewV4().String(),
				Title:       "MacBook Pro 16 M1 Max 64GB",
				Description: "Apple M1 Max chip 10-core CPU with 8 performance cores and 2 efficiency cores 32-core GPU 16-core Neural Engine 400GB/s memory bandwidth",
				Quantity:    1,
				Price:       5200,
			},
		},
	})
	if err != nil {
		appLogger.Errorf("(client.CreateOrder) err: {%v}", err)
	}
	require.NoError(t, err)
	require.NotNil(t, res)
	appLogger.Infof("result: %s", res.String())

	aggregateID = res.GetAggregateID()

	wg := &sync.WaitGroup{}
	for i := 0; i <= 60; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			result, err := client.UpdateOrder(context.Background(), &orderService.UpdateOrderReq{
				AggregateID: aggregateID,
				ShopItems: []*orderService.ShopItem{
					{
						ID:          uuid.NewV4().String(),
						Title:       "MacBook Pro 16 M1 Max 64GB",
						Description: "Apple M1 Max chip 10-core CPU with 8 performance cores and 2 efficiency cores 32-core GPU 16-core Neural Engine 400GB/s memory bandwidth",
						Quantity:    1,
						Price:       5200,
					},
					{
						ID:          uuid.NewV4().String(),
						Title:       "MacBook Pro 14 M1 Max 64GB",
						Description: "Apple M1 Max chip 10-core CPU with 8 performance cores and 2 efficiency cores 32-core GPU 16-core Neural Engine 400GB/s memory bandwidth",
						Quantity:    1,
						Price:       4200,
					},
					{
						ID:          uuid.NewV4().String(),
						Title:       "IPad PRO",
						Description: "Extreme dynamic range comes to the 12.9-inch iPad Pro.2 The Liquid Retina XDR display delivers true-to-life detail with a 1,000,000:1 contrast ratio, great for viewing and editing HDR photos and videos or enjoying your favorite movies and TV shows. It also features a breathtaking 1000 nits of full‑screen brightness and 1600 nits of peak brightness. And advanced display technologies like P3 wide color, True Tone, and ProMotion.",
						Quantity:    1,
						Price:       3200,
					},
				},
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
					ShopItems: []*orderService.ShopItem{
						{
							ID:          uuid.NewV4().String(),
							Title:       "MacBook Pro 16 M1 Max 64GB",
							Description: "Apple M1 Max chip 10-core CPU with 8 performance cores and 2 efficiency cores 32-core GPU 16-core Neural Engine 400GB/s memory bandwidth",
							Quantity:    1,
							Price:       5200,
						},
						{
							ID:          uuid.NewV4().String(),
							Title:       "MacBook Pro 14 M1 Max 64GB",
							Description: "Apple M1 Max chip 10-core CPU with 8 performance cores and 2 efficiency cores 32-core GPU 16-core Neural Engine 400GB/s memory bandwidth",
							Quantity:    1,
							Price:       4200,
						},
						{
							ID:          uuid.NewV4().String(),
							Title:       "IPad PRO",
							Description: "Extreme dynamic range comes to the 12.9-inch iPad Pro.2 The Liquid Retina XDR display delivers true-to-life detail with a 1,000,000:1 contrast ratio, great for viewing and editing HDR photos and videos or enjoying your favorite movies and TV shows. It also features a breathtaking 1000 nits of full‑screen brightness and 1600 nits of peak brightness. And advanced display technologies like P3 wide color, True Tone, and ProMotion.",
							Quantity:    1,
							Price:       3200,
						},
					},
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
				ShopItems: []*orderService.ShopItem{
					{
						ID:          uuid.NewV4().String(),
						Title:       "MacBook Pro 16 M1 Max 64GB",
						Description: "Apple M1 Max chip 10-core CPU with 8 performance cores and 2 efficiency cores 32-core GPU 16-core Neural Engine 400GB/s memory bandwidth",
						Quantity:    1,
						Price:       5200,
					},
					{
						ID:          uuid.NewV4().String(),
						Title:       "MacBook Pro 14 M1 Max 64GB",
						Description: "Apple M1 Max chip 10-core CPU with 8 performance cores and 2 efficiency cores 32-core GPU 16-core Neural Engine 400GB/s memory bandwidth",
						Quantity:    1,
						Price:       4200,
					},
					{
						ID:          uuid.NewV4().String(),
						Title:       "IPad PRO",
						Description: "Extreme dynamic range comes to the 12.9-inch iPad Pro.2 The Liquid Retina XDR display delivers true-to-life detail with a 1,000,000:1 contrast ratio, great for viewing and editing HDR photos and videos or enjoying your favorite movies and TV shows. It also features a breathtaking 1000 nits of full‑screen brightness and 1600 nits of peak brightness. And advanced display technologies like P3 wide color, True Tone, and ProMotion.",
						Quantity:    1,
						Price:       3200,
					},
				},
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
