package file

import (
	"net/http"

	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"
	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	logger := middleware.GetLogger(c)

	db := middleware.GetDB(c)
	results := []*models.File{}
	if err := db.ListFile(models.FileCollection, results); err != nil {
		logger.Error("List files failed:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, results)
}
