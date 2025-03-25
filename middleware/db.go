package middleware

import (
	"github.com/0987363/2table-backend/models"
	"github.com/gin-gonic/gin"
)

const dbKey = "Db"

var db models.DB

func ConnectDB(dataPath string) (err error) {
	db, err = models.NewBadger(dataPath)
	return err
}

func DbConnector() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(dbKey, db)
		c.Next()
	}
}

func GetDB(c *gin.Context) models.DB {
	return c.MustGet(dbKey).(models.DB)
}
