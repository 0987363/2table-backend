package user

import (
	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"net/http"
)

// Me mine user info
func MeV3(c *gin.Context) {
	db := middleware.GetDB(c)
	//	logger := middleware.GetLogger(c)

	user := models.FindUserByID(c, db, models.ObjectIDHex(middleware.GetUserID(c)))
	newUser := &models.NewUser{}
	copier.Copy(&newUser, user)
	c.JSON(http.StatusOK, newUser)
}
