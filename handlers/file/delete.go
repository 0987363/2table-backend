package file

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context) {
	logger := middleware.GetLogger(c)

	fileID := c.Param("id")

	file := models.File{}
	db := middleware.GetDB(c)
	if err := db.GetFile(models.FileCollection, fileID, &file); err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			c.Status(http.StatusNoContent)
			return
		}

		logger.Errorf("Get file by id:%s failed:%v", fileID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	storage := middleware.GetStorage(c)
	if err := storage.Delete(context.Background(), file.Path); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Errorf("Delete storage failed:%v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	if err := db.DeleteFile(models.FileCollection, fileID); err != nil {
		logger.Errorf("Delete file by id:%s failed:%v", fileID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
