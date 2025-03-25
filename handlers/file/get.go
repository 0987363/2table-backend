package file

import (
	"context"
	"mime"
	"net/http"

	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	logger := middleware.GetLogger(c)

	fileID := c.Param("id")

	file := models.File{}
	db := middleware.GetDB(c)
	if err := db.Get(models.FileCollection, fileID, &file); err != nil {
		logger.Errorf("Find file by id:%s failed:%v", fileID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	storage := middleware.GetStorage(c)
	data, err := storage.ReadAll(context.Background(), file.Path)
	if err != nil {
		logger.Errorf("Read file by path:%s failed:%s", file.Path, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	contentType := mime.TypeByExtension(file.Type)
	if contentType == "" {
		contentType = "application/octet-stream" // 默认二进制流
	}
	c.Data(http.StatusOK, contentType, data)
}
