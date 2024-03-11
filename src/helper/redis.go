package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
	"url-shortener/config"
)

type redisWrapper struct {
	Context context.Context
	Config  config.AppConfig
	Client  *redis.Client
}

func NewRedisWrapper(ctx context.Context, config config.AppConfig) IRedisWrapper {
	return &redisWrapper{
		Context: ctx,
		Config:  config,
	}
}

type IRedisWrapper interface {
	Connect() *redis.Client
	GetJson(ctx context.Context, key string) (map[string]interface{}, error)
	SetJson(ctx context.Context, key string, value interface{}, ttl int) error
}

func (r *redisWrapper) Connect() *redis.Client {
	logger := r.Context.Value("logger").(*zap.SugaredLogger)

	client := redis.NewClient(&redis.Options{
		Addr:     r.Config.RedisConfig.Uri,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	logger.Info("redis initialized successfully")
	r.Client = client
	return client
}

func (r *redisWrapper) GetJson(ctx context.Context, key string) (map[string]interface{}, error) {
	var result map[string]interface{}

	if key == "" {
		return nil, errors.New("key cannot be empty")
	}

	value := r.Client.Get(ctx, key)

	if value.Val() == "" {
		return nil, nil
	}

	err := json.Unmarshal([]byte(value.Val()), &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *redisWrapper) SetJson(ctx context.Context, key string, value interface{}, ttl int) error {
	if key == "" || value == nil {
		return errors.New("key or value cannot be empty")
	}

	bytes, err := json.Marshal(value)

	if err != nil {
		return err
	}

	status := r.Client.Set(ctx, fmt.Sprintf("%s", key), string(bytes), time.Duration(ttl))

	if status.Err() != nil {
		return err
	}

	return nil
}
