package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"url-shortener/config"
	"url-shortener/src/model"
	"url-shortener/src/service"
)

type shortenURLController struct {
	config.AppConfig
	IBaseController
	shortenURLService service.IShortenURLService
}

func NewShortenURLController(appConfig config.AppConfig, baseController IBaseController, shortenURLService service.IShortenURLService) IShortenURLController {
	return &shortenURLController{
		appConfig,
		baseController,
		shortenURLService,
	}
}

type IShortenURLController interface {
	Get(ctx *gin.Context)
	Create(ctx *gin.Context)
}

func (s shortenURLController) Get(ctx *gin.Context) {
	var response *model.Response
	var error *model.Error
	var logger *zap.SugaredLogger

	logger = ctx.Request.Context().Value("logger").(*zap.SugaredLogger)
	logger.Info("request started")

	defer func() {
		s.HandleResponse(*ctx, response, error)
	}()

	shortUrlId := ctx.Param("shortUrl")

	link, err := s.shortenURLService.Get(ctx, shortUrlId)

	if err != nil {
		error = &model.Error{
			Code:       0,
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
		return
	}

	response = &model.Response{
		StatusCode: 200,
		Code:       0,
		Data:       link,
	}
	return
}

func (s shortenURLController) Create(ctx *gin.Context) {

}
