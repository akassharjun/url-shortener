package cmd

import (
	"context"
	"errors"
	"fmt"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/integrations/nrzap"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"url-shortener/config"
	"url-shortener/src/controller"
	rw "url-shortener/src/helper"
	"url-shortener/src/log"
	"url-shortener/src/model/cache"
	"url-shortener/src/model/db"
	"url-shortener/src/router"
	"url-shortener/src/service"
)

func initLogger() *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	logFile, _ := os.OpenFile("./logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.InfoLevel

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}

func Init() {
	zapLogger := initLogger()

	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	logger = logger.With(zap.String("application", "url-shortener"))

	requiredEnvs := []string{
		"APPLICATION_NAME",
		"NEW_RELIC_ENABLED",
	}

	// load application configuration
	appConfig := config.NewAppConfig(requiredEnvs, *logger)

	r := gin.Default()

	// setup newrelic
	if appConfig.NewRelicConfig.Enabled {
		nrApp, err := newrelic.NewApplication(
			newrelic.ConfigAppName(appConfig.NewRelicConfig.AppName),
			newrelic.ConfigLicense(appConfig.NewRelicConfig.Licence),
			newrelic.ConfigAppLogForwardingEnabled(true),
			nrzap.ConfigLogger(zapLogger.Named("newrelic")),

			func(cfg *newrelic.Config) {
				cfg.CustomInsightsEvents.Enabled = false
				cfg.DistributedTracer.Enabled = true
			},
		)

		if err != nil {
			fmt.Println("unable to create New Relic Application", err)
			logger.With(err).Error("unable to create new relic application")
			return
		}

		r.Use(nrgin.Middleware(nrApp))
		r.Use(SetNewRelicInContext())
	}

	// setup middleware
	r.Use(ginzap.RecoveryWithZap(zapLogger, false))
	r.Use(UseJson)
	r.Use(SetRequestId)
	r.Use(log.NewZapGinRequestLogger(zapLogger))
	r.Use(SetLoggerInContext(logger))

	// services
	//httpClient := lib.NewHttpClient()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "logger", logger)

	redisWrapper := rw.NewRedisWrapper(ctx, appConfig)
	redisConnSuccess := redisWrapper.Connect()

	if redisConnSuccess == nil {
		err := errors.New("unable to connect to redis")
		logger.With(err).Error("unable to connect to redis")
		return
	}

	mongoWrapper := rw.NewMongoWrapper(ctx, appConfig)
	err := mongoWrapper.Connect()

	if err != nil {
		logger.With(err).Error("unable to connect to redis")
		return
	}

	database := mongoWrapper.Client.Database("url-shortener")

	if database == nil {
		err := errors.New("mongo database is missing")
		logger.With(err).Error("mongo database is missing")
		return
	}

	linkStore := db.NewLinkStore(appConfig, *database)

	linkCache := cache.NewLinkCache(appConfig, redisWrapper)

	shortenURLService := service.NewShortenURLService(appConfig, linkStore, linkCache)

	baseController := controller.NewBaseController(appConfig)
	shortenURLController := controller.NewShortenURLController(appConfig, baseController, shortenURLService)

	// setup routes
	router.StartRoutes(r, appConfig, shortenURLController)

	r.Run()
}
