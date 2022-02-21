package repository

import (
	"context"

	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/order/models"
	"github.com/AleksK1NG/es-microservice/pkg/constants"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	log logger.Logger
	cfg *config.Config
	db  *mongo.Client
}

func NewMongoRepository(log logger.Logger, cfg *config.Config, db *mongo.Client) *mongoRepository {
	return &mongoRepository{log: log, cfg: cfg, db: db}
}

func (m *mongoRepository) Insert(ctx context.Context, order *models.OrderProjection) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoRepository.Insert")
	defer span.Finish()
	span.LogFields(log.String("OrderID", order.OrderID))

	_, err := m.getOrdersCollection().InsertOne(ctx, order, &options.InsertOneOptions{})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return order.OrderID, nil
}

func (m *mongoRepository) GetByID(ctx context.Context, orderID string) (*models.OrderProjection, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoRepository.GetByID")
	defer span.Finish()
	span.LogFields(log.String("OrderID", orderID))

	var orderProjection models.OrderProjection
	if err := m.getOrdersCollection().FindOne(ctx, bson.M{constants.OrderId: orderID}).Decode(&orderProjection); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &orderProjection, nil
}

func (m *mongoRepository) UpdateOrder(ctx context.Context, order *models.OrderProjection) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoRepository.UpdateShoppingCart")
	defer span.Finish()
	span.LogFields(log.String("OrderID", order.OrderID))

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(false)

	var res models.OrderProjection
	if err := m.getOrdersCollection().FindOneAndUpdate(ctx, bson.M{constants.OrderId: order.OrderID}, bson.M{"$set": order}, ops).Decode(&res); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	m.log.Debugf("(UpdateShoppingCart) result OrderID: {%s}", res.OrderID)
	return nil
}

func (m *mongoRepository) UpdateCancel(ctx context.Context, order *models.OrderProjection) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoRepository.UpdateCancel")
	defer span.Finish()
	span.LogFields(log.String("OrderID", order.OrderID))

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(false)

	update := bson.M{"$set": bson.M{constants.Canceled: order.Canceled, constants.CancelReason: order.CancelReason}}
	var res models.OrderProjection
	if err := m.getOrdersCollection().FindOneAndUpdate(ctx, bson.M{constants.OrderId: order.OrderID}, update, ops).Decode(&res); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	m.log.Debugf("(UpdateCancel) result OrderID: {%s}", res.OrderID)
	return nil
}

func (m *mongoRepository) UpdatePayment(ctx context.Context, order *models.OrderProjection) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoRepository.UpdatePayment")
	defer span.Finish()
	span.LogFields(log.String("OrderID", order.OrderID))

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(false)

	update := bson.M{"$set": bson.M{constants.Payment: order.Payment, constants.Paid: order.Paid}}
	var res models.OrderProjection
	if err := m.getOrdersCollection().FindOneAndUpdate(ctx, bson.M{constants.OrderId: order.OrderID}, update, ops).Decode(&res); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	m.log.Debugf("(UpdatePayment) result OrderID: {%s}", res.OrderID)
	return nil
}

func (m *mongoRepository) Complete(ctx context.Context, order *models.OrderProjection) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoRepository.Complete")
	defer span.Finish()
	span.LogFields(log.String("OrderID", order.OrderID))

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(false)

	update := bson.M{"$set": bson.M{constants.Completed: order.Completed, constants.DeliveredTime: order.DeliveredTime}}
	var res models.OrderProjection
	if err := m.getOrdersCollection().FindOneAndUpdate(ctx, bson.M{constants.OrderId: order.OrderID}, update, ops).Decode(&res); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	m.log.Debugf("(Complete) result OrderID: {%s}", res.OrderID)
	return nil
}

func (m *mongoRepository) UpdateDeliveryAddress(ctx context.Context, order *models.OrderProjection) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoRepository.UpdateDeliveryAddress")
	defer span.Finish()
	span.LogFields(log.String("OrderID", order.OrderID))

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(false)

	update := bson.M{"$set": bson.M{constants.DeliveryAddress: order.DeliveryAddress}}
	var res models.OrderProjection
	if err := m.getOrdersCollection().FindOneAndUpdate(ctx, bson.M{constants.OrderId: order.OrderID}, update, ops).Decode(&res); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	m.log.Debugf("(UpdateDeliveryAddress) result OrderID: {%s}", res.OrderID)
	return nil
}

func (m *mongoRepository) UpdateSubmit(ctx context.Context, order *models.OrderProjection) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoRepository.UpdateSubmit")
	defer span.Finish()
	span.LogFields(log.String("OrderID", order.OrderID))

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(false)

	update := bson.M{"$set": bson.M{constants.Submitted: order.Submitted}}
	var res models.OrderProjection
	if err := m.getOrdersCollection().FindOneAndUpdate(ctx, bson.M{constants.OrderId: order.OrderID}, update, ops).Decode(&res); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	m.log.Debugf("(UpdateSubmit) result OrderID: {%s}", res.OrderID)
	return nil
}

func (m *mongoRepository) getOrdersCollection() *mongo.Collection {
	return m.db.Database(m.cfg.Mongo.Db).Collection(m.cfg.MongoCollections.Orders)
}
