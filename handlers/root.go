package handlers

import (
	"github.com/0987363/2table-backend/handlers/file"
	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"net/http"
	"time"
)

var RootMux = gin.New()

func Init(origins []string) {
	gin.SetMode(gin.DebugMode)

	RootMux.Use(middleware.RequestID())
	RootMux.Use(middleware.Logger())
	RootMux.Use(middleware.Recoverer())
	RootMux.Use(middleware.DbConnector())
	RootMux.Use(middleware.StorageConnector())

	if len(origins) == 0 {
		origins = []string{"*"}
	}
	RootMux.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "X-Druid-Authentication", "X-Druid-Temporary-Token", "Content-Type", "X-Result-Limit", "X-Result-Offset", "X-Result-Sort", "Accept-Language", "Accept", "If-None-Match"},
		ExposeHeaders:    []string{"X-Druid-Authentication", "X-Druid-Temporary-Token", "X-Result-Length", "X-Result-Count", "Content-Length", "Content-Type", "Content-Disposition", "Etag"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	RootMux.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, &models.BuildInfo)
	})

	v1Mux := RootMux.Group("/v1")
	{
		fileMux := v1Mux.Group("/file")
		{
			fileMux.GET("/id/:id", file.Get) // 获取某个固件信息

			fileMux.POST("/name/:name", file.Upload) // 上传固件
		}
	}
}
