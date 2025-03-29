package mineru

import (
	"context"
	"net/http"

	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Create(c *gin.Context) {
	logger := middleware.GetLogger(c)

	fileID := c.Param("id")

	file := &models.File{}
	db := middleware.GetDB(c)
	if err := db.GetFile(models.FileCollection, fileID, file); err != nil {
		logger.Errorf("Get file by id:%s failed:%v", fileID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	storage := middleware.GetStorage(c)
	if file.PDFPath == "" {
		if err := file.ConvertPDF(storage); err != nil {
			logger.Error("Convert pdf failed:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if err := db.InsertFile(models.FileCollection, file.ID, file); err != nil {
			logger.Error("Update file status failed:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	taskData := &models.TaskRequestData{
		EnableFormula: true,
		Language:      "ch",
		LayoutModel:   "doclayout_yolo",
		EnableTable:   true,
		Files: []models.FileInfo{
			{Name: file.PDFPath, IsOCR: true, DataID: file.ID},
		},
	}
	mineru := models.NewMineru(viper.GetString("ocr.mineru.token"))
	id, urls, err := mineru.CreateTask(taskData)
	if err != nil {
		logger.Error("Create task failed:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	logger.Info("Create task:", id)

	reader, err := storage.NewReader(context.Background(), file.PDFPath, nil)
	if err != nil {
		logger.Error("Create reader failed:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	if err := mineru.UploadFile(urls[0], reader); err != nil {
		logger.Error("Upload file failed:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
