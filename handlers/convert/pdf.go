package convert

import (
	"context"
	"fmt"
	"image"
	"net/http"

	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"
	"github.com/gin-gonic/gin"
	"gocloud.dev/blob"
)

func PDF(c *gin.Context) {
	logger := middleware.GetLogger(c)

	fileID := c.Param("id")

	file := models.File{}
	db := middleware.GetDB(c)
	if err := db.GetFile(models.FileCollection, fileID, &file); err != nil {
		logger.Errorf("Get file by id:%s failed:%v", fileID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if file.PDFPath != "" {
		c.Status(http.StatusNoContent)
		return
	}

	storage := middleware.GetStorage(c)
	reader, err := storage.NewReader(context.Background(), file.Path, nil)
	if err != nil {
		logger.Errorf("Init reader failed:%v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		logger.Errorf("Decode image failed:%v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	logger.Infof("Start convert:%s to pdf", file.Path)

	file.PDFPath = fmt.Sprintf("%s.pdf", file.ID)
	opts := &blob.WriterOptions{ContentType: "application/pdf"}
	writer, err := storage.NewWriter(context.Background(), file.PDFPath, opts)
	if err != nil {
		logger.Error("Init writer failed:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err := models.Image2Pdf(m, writer); err != nil {
		logger.Error("Convert pdf failed:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err := writer.Close(); err != nil {
		logger.Error("Close writer failed:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
