package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
	"url-shortener/config"
	"url-shortener/src/model/cache"
	"url-shortener/src/model/dao"
	"url-shortener/src/model/db"
)

type shortenURLService struct {
	appConfig config.AppConfig
	linkStore db.ILinkStore
	linkCache cache.ILinkCache
}

func NewShortenURLService(appConfig config.AppConfig, linkStore db.ILinkStore, linkCache cache.ILinkCache) IShortenURLService {
	return shortenURLService{
		appConfig: appConfig,
		linkStore: linkStore,
		linkCache: linkCache,
	}
}

type IShortenURLService interface {
	Get(ctx context.Context, id string) (*dao.Link, error)
}

func (s shortenURLService) Get(ctx context.Context, shortUrl string) (*dao.Link, error) {
	logger := ctx.Value("logger").(*zap.SugaredLogger)

	logger = logger.With(zap.String("shortUrl", shortUrl))

	link, err := s.linkCache.Get(ctx, shortUrl)

	if err != nil {
		error := fmt.Errorf("error while reading cache for key %s", shortUrl)
		logger.With(err).Error(error)
		return nil, error
	}

	if link != nil {
		logger.Infof("cache hit miss for key %s", shortUrl)
		return link, nil
	}

	query := dao.LinkQuery{ShortURL: &shortUrl}

	link, err = s.linkStore.FindOne(ctx, query)

	if err != nil {
		error := fmt.Errorf("error while reading mongo for short url %s", shortUrl)
		logger.With(err).Error(error)
		return nil, error
	}

	err = s.linkCache.Set(link.ShortURL, *link, int(time.Minute*100))

	if err != nil {
		logger.With(err).Error(fmt.Sprintf("error while setting cache for key %s", shortUrl))
	}

	logger.Info(fmt.Sprintf("wrote to redis cache, key: %s", shortUrl))

	return link, nil
}
