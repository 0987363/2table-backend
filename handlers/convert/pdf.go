package convert

import (
	"net/http"

	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"
	"github.com/gin-gonic/gin"
)

func PDF(c *gin.Context) {
	logger := middleware.GetLogger(c)

	fileID := c.Param("id")

	file := &models.File{}
	db := middleware.GetDB(c)
	if err := db.GetFile(models.FileCollection, fileID, file); err != nil {
		logger.Errorf("Get file by id:%s failed:%v", fileID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if file.PDFPath != "" {
		c.Status(http.StatusNoContent)
		return
	}

	storage := middleware.GetStorage(c)
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

	c.Status(http.StatusOK)
}
