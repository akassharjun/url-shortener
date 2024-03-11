package cache

import (
	"context"
	"encoding/json"
	"url-shortener/config"
	"url-shortener/src/helper"
	"url-shortener/src/model/dao"
)

type linkCache struct {
	appConfig    config.AppConfig
	redisWrapper helper.IRedisWrapper
}

func NewLinkCache(appConfig config.AppConfig, redisWrapper helper.IRedisWrapper) ILinkCache {
	return linkCache{appConfig, redisWrapper}
}

type ILinkCache interface {
	Get(ctx context.Context, key string) (*dao.Link, error)
	Set(ctx context.Context, key string, value dao.Link, expirationSeconds int) error
}

func (l linkCache) Get(ctx context.Context, key string) (*dao.Link, error) {
	result, err := l.redisWrapper.GetJson(ctx, key)

	if err != nil {
		return nil, err
	}

	var link *dao.Link
	helper.ConvertMapToStruct(result, link)

	return link, nil
}

func (l linkCache) Set(ctx context.Context, key string, value dao.Link, expirationSeconds int) error {
	jsonData, err := json.Marshal(value)

	if err != nil {
		return err
	}

	err = l.redisWrapper.SetJson(ctx, key, string(jsonData), expirationSeconds)

	if err != nil {
		return err
	}

	return nil
}
