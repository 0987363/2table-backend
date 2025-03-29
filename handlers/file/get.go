package file

import (
	"context"
	"fmt"
	"io"
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
	if err := db.GetFile(models.FileCollection, fileID, &file); err != nil {
		logger.Errorf("Get file by id:%s failed:%v", fileID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	storage := middleware.GetStorage(c)
	reader, err := storage.NewReader(context.Background(), file.ID, nil)
	if err != nil {
		logger.Errorf("Init reader failed:%v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	c.Header("Content-Type", reader.ContentType())
	c.Header("Content-Length", fmt.Sprintf("%d", file.Size))
	c.Header("Accept-Ranges", "bytes")

	if _, err = io.CopyBuffer(c.Writer, reader, make([]byte, 64<<10)); err != nil {
		logger.Error("Download failed:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
