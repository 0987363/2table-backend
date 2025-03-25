package middleware

import (
	"context"

	"github.com/0987363/2table-backend/models"
	"github.com/gin-gonic/gin"
)

const storageKey = "Storage"

var sm *models.StorageManager

func ConnectStorageManager(cfg *models.StorageConfig) (err error) {
	sm, err = models.NewStorageManager(context.Background(), cfg)
	return nil
}

func StorageConnector() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(storageKey, sm)

		c.Next()
	}
}

func GetStorage(c *gin.Context) *models.StorageManager {
	return c.MustGet(storageKey).(*models.StorageManager)
}
