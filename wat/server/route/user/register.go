package wat

import (
	"github.com/ChikyuKido/wat/wat/helper"
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	util "github.com/ChikyuKido/wat/wat/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerData = struct {
			Email    string `json:"email"`
			Username string `json:"username"`
			Password string `json:"password"`
		}{}
		if err := c.ShouldBind(&registerData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request data"})
			return
		}
		if registerData.Email == "guest" || registerData.Username == "guest" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "guest is a reserved user"})
			return
		}
		if strings.TrimSpace(registerData.Email) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email must not be empty"})
			return
		}
		if !helper.IsValidEmail(registerData.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
			return
		}
		found := false
		for _, email := range util.Config.AllowedEmails {
			if email == strings.Split(registerData.Email, "@")[1] {
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": "the email domain is not allowed"})
			return
		}

		if strings.TrimSpace(registerData.Username) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username must not be empty"})
			return
		}
		if strings.TrimSpace(registerData.Password) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must not be empty"})
			return
		}
		if repo.DoesUserByEmailExist(registerData.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exist"})
			return
		}
		if repo.DoesUserByUsernameExist(registerData.Username) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exist"})
			return
		}
		hashedPassword, err := util.HashPassword(registerData.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			logrus.Errorf("Failed to hash password: %v", err)
			return
		}
		if !repo.InsertNewUser(registerData.Username, hashedPassword, registerData.Email) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create new user"})
			logrus.Errorf("Failed to create new user: %v", err)
			return
		}
		if !util.Config.FirstUser {
			user := repo.GetUserByEmail(registerData.Email)
			if util.Config.EmailVerification {
				emailSend := helper.SendEmailVerificationForUser(user)
				c.JSON(http.StatusOK, gin.H{"message": "successful create an account", "verification": true, "emailSent": emailSend})
				return
			} else {
				if !repo.AddRoleToUser(user.ID, repo.GetRoleByName("user").ID) { // user
					logrus.Errorf("Failed to assign roles to a newly created user. User now has zero permissions")
				}
				repo.VerifyUser(user.ID)
			}
		} else {
			user := repo.GetUserByEmail(registerData.Email)
			if !repo.AddRoleToUser(user.ID, repo.GetRoleByName("admin").ID) { // admin
				logrus.Errorf("Failed to assign roles to a newly created user. User now has zero permissions")
			}
			util.Config.FirstUser = false
		}
		c.JSON(http.StatusOK, gin.H{"message": "successful create an account", "verification": false})
		return
	}
}
