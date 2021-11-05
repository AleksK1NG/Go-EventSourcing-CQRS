package repository

import (
	"context"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/models"
	"github.com/AleksK1NG/es-microservice/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository interface {
	Insert(ctx context.Context, order *models.OrderProjection) (string, error)
	GetByID(ctx context.Context, orderID string) (*models.OrderProjection, error)
	UpdateOrder(ctx context.Context, order *models.OrderProjection) error
}

type mongoRepository struct {
	log logger.Logger
	cfg *config.Config
	db  *mongo.Client
}

func NewMongoRepository(log logger.Logger, cfg *config.Config, db *mongo.Client) *mongoRepository {
	return &mongoRepository{log: log, cfg: cfg, db: db}
}

func (m *mongoRepository) Insert(ctx context.Context, order *models.OrderProjection) (string, error) {
	collection := m.db.Database(m.cfg.Mongo.Db).Collection(m.cfg.MongoCollections.Orders)

	_, err := collection.InsertOne(ctx, order, &options.InsertOneOptions{})
	if err != nil {
		return "", err
	}

	return order.OrderID, nil
}

func (m *mongoRepository) GetByID(ctx context.Context, orderID string) (*models.OrderProjection, error) {
	collection := m.db.Database(m.cfg.Mongo.Db).Collection(m.cfg.MongoCollections.Orders)

	var orderProjection models.OrderProjection
	if err := collection.FindOne(ctx, bson.M{"orderId": orderID}).Decode(&orderProjection); err != nil {
		return nil, err
	}

	return &orderProjection, nil
}

func (m *mongoRepository) UpdateOrder(ctx context.Context, order *models.OrderProjection) error {
	collection := m.db.Database(m.cfg.Mongo.Db).Collection(m.cfg.MongoCollections.Orders)

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(false)

	res := make(map[string]interface{}, 10)
	if err := collection.FindOneAndUpdate(ctx, bson.M{"orderId": order.OrderID}, bson.M{"$set": order}, ops).Decode(&res); err != nil {
		return err
	}

	m.log.Debugf("update result: %+v", res)
	return nil
}
