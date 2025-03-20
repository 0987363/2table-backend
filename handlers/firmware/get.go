package firmware

import (
	"net/http"

	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"
	"github.com/gin-gonic/gin"
)

func getFirmwareValidate(c *gin.Context) string {
	firmwareID := c.Param("id")
	if !models.IsObjectIDHex(firmwareID) {
		c.AbortWithError(http.StatusBadRequest, models.Error("Firmware id is invalid.", firmwareID))
		return ""
	}

	return firmwareID
}

// Get all of firmware file
func Get(c *gin.Context) {
	db := middleware.GetDB(c)
	logger := middleware.GetLogger(c)

	firmwareID := getFirmwareValidate(c)
	if firmwareID == "" {
		return
	}

	user := models.FindUserByID(c, db, models.ObjectIDHex(middleware.GetUserID(c)))
	if user.IsManager() && !user.CheckFirmwareRead() {
		logger.Error("You did not have permissions read firmware.")
		c.AbortWithError(http.StatusForbidden, models.Error("You did not have permissions to read firmware."))
		return
	}

	firmware, err := models.FindFirmwareByID(c, db, models.ObjectIDHex(firmwareID))
	if err != nil {
		logger.Error("Find firmware failed by id:", firmwareID, err)
		c.AbortWithError(http.StatusInternalServerError, models.Error("Find firmware failed by id:", firmwareID, err))
		return
	}

	c.JSON(http.StatusOK, firmware)
}
