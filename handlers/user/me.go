package user

import (
	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Me mine user info
func Me(c *gin.Context) {
	db := middleware.GetDB(c)
	//	logger := middleware.GetLogger(c)
	user := models.FindUserByID(c, db, models.ObjectIDHex(middleware.GetUserID(c)))
	if user.IsManager() {
		user.Permissions.ManagerPermissions = nil
	}
	c.JSON(http.StatusOK, user)
}
