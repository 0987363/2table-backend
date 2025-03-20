package handlers

import (
	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"

	"github.com/0987363/2table-backend/handlers/firmware"
	"github.com/0987363/2table-backend/handlers/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"net/http"
	"time"
)

var RootMux = gin.New()

func Init(origins []string) {
	gin.SetMode(gin.DebugMode)

	RootMux.Use(middleware.ClientInfo())
	RootMux.Use(middleware.RequestID())
	RootMux.Use(middleware.Logger())
	RootMux.Use(middleware.Recoverer())
	RootMux.Use(middleware.DBConnector())
	RootMux.Use(middleware.Redis())
	RootMux.Use(middleware.Micro())
	RootMux.Use(middleware.StackClient())

	/*
		RootMux.Use(cors.Middleware(cors.Config{
			//		Origins:         "*",
			Origins:         "https://bird.coolhei.com, http://bird.coolhei.com, https://bird.druidtech.cn, https://bird.druid.tech",
			Methods:         "GET, PUT, POST, DELETE",
			RequestHeaders:  "Origin, X-Druid-Authentication, Content-Type, X-Result-Limit, X-Result-Offset, X-Result-Sort, Accept-Language, Accept",
			ExposedHeaders:  "X-Druid-Authentication, X-Result-Count, Content-Length, Content-Type, Content-Disposition",
			MaxAge:          time.Hour * 24 * 7,
			Credentials:     true,
			ValidateHeaders: false,
		}))
	*/

	if len(origins) == 0 {
		origins = []string{"*"}
	}
	RootMux.Use(cors.New(cors.Config{
		//			AllowOrigins:     []string{"https://bird.coolhei.com", "http://bird.coolhei.com", "https://bird.druidtech.cn", "https://bird.druid.tech"},
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "X-Druid-Authentication", "X-Druid-Temporary-Token", "Content-Type", "X-Result-Limit", "X-Result-Offset", "X-Result-Sort", "Accept-Language", "Accept", "If-None-Match"},
		ExposeHeaders:    []string{"X-Druid-Authentication", "X-Druid-Temporary-Token", "X-Result-Length", "X-Result-Count", "Content-Length", "Content-Type", "Content-Disposition", "Etag"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	/*
		cache.SetPageKey("manager:1:" + viper.GetString("release"))
		var store persistence.CacheStore
		memAddrs := viper.GetString("memcached")
		addrs := strings.Split(memAddrs, ",")
		if memAddrs == "" || len(addrs) == 0 {
			store = persistence.NewInMemoryStore(time.Hour * 12)
		} else {
			memStore := persistence.NewMemcachedStore(addrs, time.Hour*12)
			memStore.Client.Timeout = time.Second * 5
			store = memStore
		}
	*/

	RootMux.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, &models.BuildInfo)
	})

	v2Mux := RootMux.Group("/v2")
	{
		{
			v2Mux.POST("/login", user.Login)
		}

		v2Mux.Use(middleware.Authenticator())

		firmwareMux := v2Mux.Group("/firmware")
		{
			firmwareMux.GET("/id/:id", firmware.Get) // 获取某个固件信息

			firmwareMux.DELETE("/id/:id", firmware.DeleteFromMongo) // 删除固件

			firmwareMux.POST("/name/:name", firmware.UploadToMongo) // 上传固件
		}
	}
}
