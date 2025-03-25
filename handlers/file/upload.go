package file

import (
	"context"
	"mime/multipart"

	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"

	"github.com/gin-gonic/gin"

	"io/ioutil"
	"net/http"

	"strings"
)

func uploadValidate(c *gin.Context) *multipart.FileHeader {
	f, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(http.StatusForbidden, models.Error("Read data from form failed.", err))
		return nil
	}

	return f
}

func Upload(c *gin.Context) {
	logger := middleware.GetLogger(c)

	f := uploadValidate(c)
	if f == nil {
		return
	}

	fileName := strings.TrimSpace(f.Filename)
	if fileName == "" {
		logger.Error("Clould not load file name.")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fd, _ := f.Open()
	data, _ := ioutil.ReadAll(fd)
	defer fd.Close()

	file := models.NewFile(fileName, data)

	db := middleware.GetDB(c)
	if err := db.Insert(models.FileCollection, file.ID, file); err != nil {
		logger.Error("Save file meta failed.", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	storage := middleware.GetStorage(c)
	if err := storage.WriteAll(context.Background(), file.Path, data, nil); err != nil {
		logger.Error("Save file failed.", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, file)
}
