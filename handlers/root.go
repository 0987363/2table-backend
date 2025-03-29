package handlers

import (
	"github.com/0987363/2table-backend/handlers/convert"
	"github.com/0987363/2table-backend/handlers/file"
	"github.com/0987363/2table-backend/handlers/ocr/mineru"
	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"

	"github.com/gin-gonic/gin"

	"net/http"
)

var RootMux = gin.New()

func init() {
	gin.SetMode(gin.DebugMode)

	RootMux.Use(middleware.RequestID())
	RootMux.Use(middleware.Logger())
	RootMux.Use(middleware.Recoverer())
	RootMux.Use(middleware.DbConnector())
	RootMux.Use(middleware.StorageConnector())

	RootMux.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, &models.BuildInfo)
	})

	v1Mux := RootMux.Group("/v1")
	{
		fileMux := v1Mux.Group("/file")
		{
			fileMux.GET("/id/:id", file.Get)
			fileMux.GET("/", file.List)

			fileMux.POST("/", file.Upload)

			fileMux.DELETE("/id/:id", file.Delete)
		}

		convertMux := v1Mux.Group("/convert")
		{
			convertMux.PUT("/file/:id/pdf", convert.PDF)
		}

		mineruMux := v1Mux.Group("/mineru")
		{
			mineruMux.POST("/file/:id", mineru.Create)
		}
	}
}
