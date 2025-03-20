package firmware

import (
	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"

	"io/ioutil"
	"net/http"

	//	"strconv"
	"strings"
)

func uploadValidate(c *gin.Context) *models.Firmware {
	firmware := models.Firmware{}
	f, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(http.StatusForbidden, models.Error("Read data from form failed.", err))
		return nil
	}

	fileName := strings.TrimSpace(f.Filename)
	if fileName == "" {
		c.AbortWithError(http.StatusForbidden, models.Error("File name is invalid."))
		return nil
	}
	firmware.Name = fileName
	firmware.File = f

	return &firmware
}

// Upload firmware
func Upload(c *gin.Context) {
	db := middleware.GetDB(c)
	logger := middleware.GetLogger(c)

	firmware := uploadValidate(c)
	if firmware == nil {
		return
	}

	user := models.FindUserByID(c, db, models.ObjectIDHex(middleware.GetUserID(c)))
	if user.IsManager() && !user.CheckFirmwareWrite() {
		logger.Error("You did not have permissions edit firmware.")
		c.AbortWithError(http.StatusForbidden, models.Error("You did not have permissions to edit firmware."))
		return
	}

	fd, err := firmware.File.Open()
	if err != nil {
		logger.Error("Open firmware file failed :", err)
		c.AbortWithError(http.StatusInternalServerError, models.Error("Open firmware file failed:", err))
		return
	}
	defer fd.Close()

	firmware.Data, err = ioutil.ReadAll(fd)
	if err != nil {
		logger.Error("Read firmware file failed :", err)
		c.AbortWithError(http.StatusInternalServerError, models.Error("Read firmware file failed:", err))
		return
	}
	firmware.Size = len(firmware.Data)

	if err := firmware.ParseFirmware(firmware.Data); err != nil {
		logger.Error("Parse firmware failed :", err)
		c.AbortWithError(http.StatusInternalServerError, models.Error("Parse firmware failed:", err))
		return
	}
	fv, err := firmware.FindFirmwareByVersion(c, db)
	if err != nil {
		logger.Error("Find firmware by version failed :", err)
		c.AbortWithError(http.StatusInternalServerError, models.Error("Find firmware by version failed:", err))
		return
	}
	if fv != nil {
		logger.Error("Had exists firmware version.")
		c.AbortWithError(http.StatusInternalServerError, models.Error("Firmware had exists version."))
		return
	}

	file := &models.File{
		ID:         primitive.NewObjectID(),
		Data:       firmware.Data,
		Size:       firmware.Size,
		Name:       firmware.Name,
		Collection: models.FilesCollection,
	}
	if err := file.Upload(c, db); err != nil {
		logger.Error("Upload file failed.", err)
		c.AbortWithError(http.StatusInternalServerError, models.Error("Upload file failed.", err))
		return
	}

	firmware.Point = file.ID
	//	firmware.MD5 = file.MD5

	if err = firmware.Upload2(c, db); err != nil {
		logger.Error("Upload firmware failed.", err)
		c.AbortWithError(http.StatusInternalServerError, models.Error("Upload firmware failed.", err))
		return
	}

	c.JSON(http.StatusCreated, firmware)
}
