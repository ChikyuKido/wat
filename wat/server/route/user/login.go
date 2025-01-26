package wat

import (
	entity "github.com/ChikyuKido/wat/wat/server/db/entity"
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	util "github.com/ChikyuKido/wat/wat/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData = struct {
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
		if err := c.ShouldBind(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request data"})
			return
		}
		if loginData.Email == nil && loginData.Username == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request data"})
			return
		}
		if (loginData.Email != nil && *loginData.Email == "guest") || (loginData.Username != nil && *loginData.Username == "guest") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "guest is a reserved user"})
			return
		}
		var user *entity.User = nil
		if loginData.Email != nil {
			if !repo.DoesUserByEmailExist(*loginData.Email) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "wrong credentials"})
				return
			}
			user = repo.GetUserByEmail(*loginData.Email)
		}
		if user == nil {
			if !repo.DoesUserByUsernameExist(*loginData.Username) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "wrong credentials"})
				return
			}
			user = repo.GetUserByUsername(*loginData.Username)
		}
		if !util.CheckPassword(user.Password, loginData.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "wrong credentials"})
			return
		}
		token, err := util.GenerateJWT(user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}
		if !user.Verified && util.Config.EmailVerification {
			id := uuid.New()
			verificationUUIDS[id.String()] = user.ID
			c.JSON(http.StatusBadRequest, gin.H{"error": "This user is not verified yet.", "verificationUUID": id.String()})
		} else {
			c.SetCookie("jwt", token, 60*60*24*30, "/", os.Getenv("DOMAIN"), false, true)
			c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in"})
		}
		return
	}
}
