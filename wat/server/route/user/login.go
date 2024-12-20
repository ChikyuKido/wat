package wat

import (
	entity "github.com/ChikyuKido/wat/wat/server/db/entity"
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	util "github.com/ChikyuKido/wat/wat/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerData = struct {
			Email    *string `json:"email"`
			Username *string `json:"username"`
			Password string  `json:"password"`
		}{}

		loggedInUser := util.GetUserFromContext(c)
		if loggedInUser == nil {
			return
		}
		if loggedInUser.Email != "guest" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already logged in"})
			return
		}
		if err := c.ShouldBind(&registerData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request data"})
			return
		}
		if registerData.Email == nil && registerData.Username == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request data"})
			return
		}
		if (registerData.Email != nil && *registerData.Email == "guest") || (registerData.Username != nil && *registerData.Username == "guest") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "guest is a reserved user"})
			return
		}
		var user *entity.User = nil
		if registerData.Email != nil {
			if !repo.DoesUserByEmailExist(*registerData.Email) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "wrong credentials"})
				return
			}
			user = repo.GetUserByEmail(*registerData.Email)
		}
		if user == nil {
			if !repo.DoesUserByUsernameExist(*registerData.Username) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "wrong credentials"})
				return
			}
			user = repo.GetUserByUsername(*registerData.Username)
		}
		if !util.CheckPassword(user.Password, registerData.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "wrong credentials"})
			return
		}
		token, err := util.GenerateJWT(*registerData.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}
		c.SetCookie("jwt", token, 60*60*24*30, "/", os.Getenv("DOMAIN"), false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in"})
		return
	}
}
