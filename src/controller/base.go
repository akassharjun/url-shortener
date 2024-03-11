package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"url-shortener/config"
	"url-shortener/src/model"
)

type baseController struct {
	config.AppConfig
}

func NewBaseController(appConfig config.AppConfig) IBaseController {
	return baseController{
		appConfig,
	}
}

type IBaseController interface {
	HandleResponse(c gin.Context, response *model.Response, error *model.Error)
}

type HTTPResponse struct {
	Success  bool            `json:"success"`
	Response *model.Response `json:"response"`
	Error    *model.Error    `json:"error"`
}

func (b baseController) HandleResponse(c gin.Context, response *model.Response, error *model.Error) {
	reqResponse := HTTPResponse{
		Success:  false,
		Response: nil,
		Error:    nil,
	}
	httpStatus := http.StatusBadGateway

	// log
	baseLogger := c.Request.Context().Value("logger").(*zap.SugaredLogger)

	defer func() {
		baseLogger.Info("request completed")
		c.JSON(httpStatus, reqResponse)
	}()

	// transform error
	if error != nil {
		reqResponse.Success = false
		reqResponse.Error = error
		httpStatus = error.StatusCode

		bytes, _ := json.Marshal(error)
		baseLogger.Error(string(bytes))

		return
	}

	// transform response
	if response != nil {
		reqResponse.Response = response
		httpStatus = response.StatusCode
	}
}
