package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"url-shortener/config"
	"url-shortener/src/controller"
)

func StartRoutes(r *gin.Engine, appConfig config.AppConfig, shortenURLController controller.IShortenURLController) {

	api := r.Group("/v1")
	{
		shortenURLRoutes := api.Group("/")
		{
			shortenURLRoutes.GET("/:shortUrl", shortenURLController.Get)
			shortenURLRoutes.POST("/create", shortenURLController.Create)
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
