package config

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"url-shortener/src/log/constants"
)

type AppConfig struct {
	RedisConfig    RedisConfig
	NewRelicConfig NewRelicConfig
	MongoConfig    MongoConfig
}

func NewAppConfig(requiredENVs []string, logger zap.SugaredLogger) AppConfig {
	baseLogger := logger.With(zap.Any("origin", constants.ApplicationInit))

	viper.AddConfigPath("url-shortener/etc")
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		baseLogger.Warn("unable to read .env file", zap.Any("error", err))
	}

	viper.SetConfigFile(".json")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		baseLogger.Warn("unable to read .json file", zap.Any("error", err))
	}

	// Get ENVIRONMENT from environment variables
	viper.AutomaticEnv()

	// Check if missing required variables
	for _, env := range requiredENVs {
		if viper.GetString(env) == "" {
			baseLogger.Error("ENV variable was not set")
			panic(fmt.Sprintf("service can't operate without %s ENV variable", env))
		}
	}

	logger.Infow("loaded ENV variables successfully")

	return AppConfig{
		RedisConfig: RedisConfig{
			Uri: viper.GetString("REDIS_URI"),
		},
		NewRelicConfig: NewRelicConfig{
			Licence: viper.GetString("NEW_RELIC_LICENSE_KEY"),
			AppName: viper.GetString("NEW_RELIC_APP_NAME"),
			Enabled: viper.GetBool("NEW_RELIC_ENABLED"),
		},
		MongoConfig: MongoConfig{
			Uri: viper.GetString("mongo.uri"),
		},
	}
}
