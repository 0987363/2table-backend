package file

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
)

func uploadValidate(c *gin.Context) (*multipart.Reader, error) {
	reader, err := c.Request.MultipartReader()
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func Upload(c *gin.Context) {
	logger := middleware.GetLogger(c)

	reader, err := uploadValidate(c)
	if err != nil {
		logger.Error("Parse multipart failed:", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	storage := middleware.GetStorage(c)
	db := middleware.GetDB(c)

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if part.FileName() == "" {
			continue
		}

		file, err := models.NewFile(db, part.FileName())
		if err != nil {
			logger.Error("Init file failed:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		opts := &blob.WriterOptions{ContentType: "application/octet-stream"}
		writer, err := storage.NewWriter(context.Background(), file.Path, opts)
		if err != nil {
			logger.Error("Init writer failed:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		bytesCopied, err := io.CopyBuffer(writer, part, make([]byte, 64<<10))
		if err != nil {
			writer.Close()

			logger.Error("Upload failed:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if err := writer.Close(); err != nil {
			logger.Error("Close writer failed:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		file.Size = bytesCopied
		file.Status = models.FileStatusFinished
		if err := db.Insert(models.FileCollection, file.ID, file); err != nil {
			logger.Error("Update file status failed:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusCreated, file)
		return
	}

	c.AbortWithStatus(http.StatusBadRequest)
}
