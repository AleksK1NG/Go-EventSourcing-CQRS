package server

import (
	"context"
	"github.com/AleksK1NG/es-microservice/pkg/constants"
	serviceErrors "github.com/AleksK1NG/es-microservice/pkg/service_errors"
	"github.com/AleksK1NG/es-microservice/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *server) initMongoDBCollections(ctx context.Context) {
	err := s.mongoClient.Database(s.cfg.Mongo.Db).CreateCollection(ctx, s.cfg.MongoCollections.Orders)
	if err != nil {
		if !utils.CheckErrMessages(err, serviceErrors.ErrMsgMongoCollectionAlreadyExists) {
			s.log.Warnf("(CreateCollection) err: {%v}", err)
		}
	}

	indexOptions := options.Index().SetSparse(true).SetUnique(true)
	index, err := s.mongoClient.Database(s.cfg.Mongo.Db).Collection(s.cfg.MongoCollections.Orders).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: constants.OrderIdIndex, Value: 1}},
		Options: indexOptions,
	})
	if err != nil {
		s.log.Warnf("(CreateOne) err: {%v}", err)
	}
	s.log.Infof("(CreatedIndex) index: {%s}", index)

	list, err := s.mongoClient.Database(s.cfg.Mongo.Db).Collection(s.cfg.MongoCollections.Orders).Indexes().List(ctx)
	if err != nil {
		s.log.Warnf("(List) err: {%v}", err)
	}

	if list != nil {
		var results []bson.M
		if err := list.All(ctx, &results); err != nil {
			s.log.Warnf("(All) err: {%v}", err)
		}
		s.log.Infof("(indexes) results: {%#v}", results)
	}

	collections, err := s.mongoClient.Database(s.cfg.Mongo.Db).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		s.log.Warnf("(ListCollections) err: {%v}", err)
	}
	s.log.Infof("(Collections) created collections: {%v}", collections)
}
