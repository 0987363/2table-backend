package firmware

import (
	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"

	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func uploadToMongoValidate(c *gin.Context) (*models.Firmware, *models.ErrorBody) {
	firmware := models.Firmware{}

	fileName := c.Param("name")
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		return nil, models.NewErrorBodyWithSub(models.CodeFirmware, []int{models.CodeFirmwareUploadNameError}, "File name is invalid.")
	}
	firmware.Name = fileName

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return nil, models.NewErrorBodyWithSub(models.CodeFirmware, []int{models.CodeFirmwareUploadFileError}, models.Error("Read data from request failed.", err))
	}

	size, err := strconv.Atoi(c.GetHeader("Content-Length"))
	if err != nil {
		return nil, models.NewErrorBodyWithSub(models.CodeCommon, []int{models.CodeCommandParameterError}, models.Error("Recv file length failed.", err))

	}
	if size == 0 || len(data) != size {
		return nil, models.NewErrorBodyWithSub(models.CodeFirmware, []int{models.CodeFirmwareUploadSizeError}, models.Error("File length invalid.", err))
	}
	firmware.Size = size
	firmware.Data = data

	return &firmware, nil
}

// UploadFirmware2 firmware
func UploadToMongo(c *gin.Context) {
	db := middleware.GetDB(c)
	logger := middleware.GetLogger(c)

	firmware, errBody := uploadToMongoValidate(c)
	if errBody != nil {
		logger.Error("Validate failed: ", errBody)
		c.AbortWithStatusJSON(http.StatusBadRequest, errBody)
		return
	}

	user := models.FindUserByID(c, db, models.ObjectIDHex(middleware.GetUserID(c)))
	if user.IsManager() && !user.CheckFirmwareWrite() {
		logger.Error("You did not have permissions edit firmware.")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if err := firmware.ParseFirmware(firmware.Data); err != nil {
		logger.Error("Parse firmware failed: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	m := bson.M{
		"device_type":      firmware.DeviceType,
		"hardware_version": firmware.HardwareVersion,
		"firmware_version": firmware.FirmwareVersion,
	}
	if firmware.FirmwareVersion >= 1000 && firmware.CodeHash != models.FirmwareCodeHashInvalidValue && firmware.CodeHash != 0 {
		//检测code_hash uint32转16进制
		codeHash := fmt.Sprintf("%x", firmware.CodeHash)
		if len(codeHash) > models.FirmwareCodeHashHexLengthMax {
			logger.Error("Code hash is invalid.", codeHash)

			c.AbortWithStatusJSON(http.StatusBadRequest, models.NewErrorBodyWithSub(models.CodeFirmware, []int{models.CodeFirmwareUploadCodeHashError}, models.Error("Code hash is invalid.")))
			return
		}
		m["build_time"] = firmware.BuildTime
		m["code_hash"] = firmware.CodeHash
	}

	fv, err := models.CheckFirmwareByFilter(db, m)
	if err != nil {
		logger.Error("Find firmware by version failed :", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if fv != nil {
		logger.Error("Had exists firmware version.")
		c.AbortWithStatusJSON(http.StatusBadRequest, models.NewErrorBodyWithSub(models.CodeFirmware, []int{models.CodeFirmwareUploadExistError}, models.Error("Upload firmware is exists.")))
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
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	firmware.Point = file.ID
	//	firmware.MD5 = file.MD5

	if err = firmware.Upload2(c, db); err != nil {
		logger.Error("Upload firmware failed.", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, firmware)
}
