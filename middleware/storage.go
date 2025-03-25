package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/0987363/2table-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const storageKey = "Storage"

var bucketURL string

func ConnectStorageManager() error {
	Type := viper.GetString("storage.type")

	switch Type {
	case models.StorageTypeLocal:
		logrus.Info("Storage Local.")

		path := viper.GetString("storage.local_config.path")
		os.MkdirAll(path, 0755)

		temp := viper.GetString("storage.local_config.temp")
		os.MkdirAll(temp, 0755)
		os.Setenv("TMPDIR", temp)

		bucketURL = "file://" + path
	case models.StorageTypeS3:
		logrus.Info("Storage S3.")
		s3 := &models.S3Config{
			Bucket:   viper.GetString("storage.local_config.path"),
			Region:   viper.GetString("storage.local_config.region"),
			Endpoint: viper.GetString("storage.local_config.endpoint"),
		}

		bucketURL = models.BuildS3URL(s3)
	default:
		return models.Error("Unknown storage type:", Type)
	}

	if viper.GetBool("storage.skip_attr") {
		logrus.Info("Storage ignore attr.")
		bucketURL += "?metadata=skip"
	}

	return nil
}

func StorageConnector() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := GetLogger(c)
		sm, err := models.NewStorageManager(context.Background(), bucketURL)
		if err != nil {
			logger.Error("Init storage failed:", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		defer sm.Close()

		c.Set(storageKey, sm)
		c.Next()
	}
}

func GetStorage(c *gin.Context) *models.StorageManager {
	return c.MustGet(storageKey).(*models.StorageManager)
}
