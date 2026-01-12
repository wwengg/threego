package store

import (
	"context"
	"time"

	"github.com/wwengg/threego/core/slog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	*mongo.Client
}

var mongoClientInstance *MongoDB

func MongoIns() *MongoDB {
	if mongoClientInstance == nil {
		slog.Ins().Errorf("mongo client is nil")
		return mongoClientInstance
	}
	return mongoClientInstance
}

func NewMongoClient(mongoURI string) (*MongoDB, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI).SetConnectTimeout(5*time.Second))
	if err != nil {
		slog.Ins().Errorf("mongo connect err %v", err)
		return nil, err
	}
	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		slog.Ins().Errorf("mongo ping err %v", err)
		return nil, err
	}
	mongoClientInstance = &MongoDB{client}
	return &MongoDB{client}, nil
}

func CloseMongoClient() {
	if err := mongoClientInstance.Disconnect(context.Background()); err != nil {
		slog.Ins().Errorf("close mongo client err %v", err)
	}
}
