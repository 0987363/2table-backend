package user

import (
	"net/http"
	//"regexp"
	"strings"
	//	"time"

	"github.com/spf13/viper"

	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"

	"github.com/gin-gonic/gin"
)

func loginValidate(c *gin.Context) *models.User {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, models.Error("Unable to parse and decode the request.", err))
		return nil
	}

	if !models.RegexpUserNameLogin.MatchString(user.UserName) {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.NewErrorBody(models.CodeUserNameInvalid, "username must be a valid value."))
		return nil
	}

	user.Password = strings.TrimSpace(user.Password)
	if !models.RegexpPwd.MatchString(user.Password) {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.NewErrorBody(models.CodePasswordInvalid, "password must be a valid value."))
		return nil
	}

	return &user
}

// Login in with username and sha256 password
func Login(c *gin.Context) {
	db := middleware.GetDB(c)
	logger := middleware.GetLogger(c)

	userReq := loginValidate(c)
	if userReq == nil {
		return
	}

	user, err := models.CheckUserByUserName(c, db, userReq.UserName)
	if err != nil {
		logger.Warn("Find user by user name failed.", err)
		c.AbortWithError(http.StatusInternalServerError, models.Error("Find user by user name failed.", err))
		return
	}
	if user == nil {
		logger.Warn("The user was invalid.")
		c.AbortWithStatusJSON(http.StatusBadRequest, models.NewErrorBody(models.CodePasswordError, ""))
		return
	}

	if user.Role != models.RoleRoot && user.Role != models.RoleManager {
		logger.Error("You don't have permission to access login.")
		c.AbortWithStatusJSON(http.StatusForbidden, models.NewErrorBody(models.CodeUserRoleForbidLogin, ""))
		return
	}

	if err := models.PasswordVerify(userReq.Password, user.HashedPassword); err != nil {
		logger.Error("Password is invalid")
		c.AbortWithStatusJSON(http.StatusBadRequest, models.NewErrorBody(models.CodePasswordError, ""))
		return
	}

	if !user.CheckTokenValid() {
		token := models.NewToken(user.ID, viper.GetString("authentication.secret"))
		user.Token = token.String()
		user.Expiry = &token.Expiry

		if err = user.UpdateToken(c, db); err != nil {
			logger.Error("Update token failed.")
			c.AbortWithError(http.StatusInternalServerError, models.Error("Update token failed."))
			return
		}
		logger.Infof("User:%s update new token in mongo.", user.UserName)
	}

	c.Header(middleware.AuthenticationHeader, user.Token)
	c.JSON(http.StatusOK, user)
}
