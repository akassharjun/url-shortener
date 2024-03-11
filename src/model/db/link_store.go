package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"url-shortener/config"
	"url-shortener/src/constants"
	"url-shortener/src/model/dao"
)

type linkStore struct {
	appConfig  config.AppConfig
	collection *mongo.Collection
}

func NewLinkStore(appConfig config.AppConfig, database mongo.Database) ILinkStore {
	linkCollection := database.Collection(constants.LinkCollection)

	return linkStore{appConfig, linkCollection}
}

type ILinkStore interface {
	FindOne(ctx context.Context, query dao.LinkQuery, options ...*options.FindOneOptions) (*dao.Link, error)
}

func (l linkStore) FindOne(ctx context.Context, query dao.LinkQuery, options ...*options.FindOneOptions) (*dao.Link, error) {
	var link dao.Link

	result := l.collection.FindOne(ctx, query, options...)

	if result.Err() != nil {
		// Check if the error is due to no documents found
		if result.Err() == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, result.Err()
	}

	if err := result.Decode(&link); err != nil {
		return nil, err
	}

	return &link, nil
}
