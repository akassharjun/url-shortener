package helper

import (
	"context"
	"go.uber.org/zap"
	"time"
	"url-shortener/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoWrapper struct {
	Context context.Context
	Config  config.AppConfig
	Client  *mongo.Client
}

func NewMongoWrapper(ctx context.Context, config config.AppConfig) *MongoWrapper {
	return &MongoWrapper{
		Context: ctx,
		Config:  config,
	}
}

type MongoFunctions interface {
	Connect()
}

func (m *MongoWrapper) Connect() error {
	logger := m.Context.Value("logger").(*zap.SugaredLogger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoConfig := m.Config.MongoConfig
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConfig.Uri))

	if err != nil {
		logger.With(err).Error("there was an error in initializing mongodb")
		return err
	}

	logger.Info("mongo initialized successfully")
	m.Client = client

	return nil
}
