package wat

import (
	util "github.com/ChikyuKido/wat/wat/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response = struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			ID       uint   `json:"id"`
			Verified bool   `json:"verified"`
		}{}
		user := util.GetUserFromContext(c)
		if user.Email == "guest" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't get the guest profile"})
		}
		response.Username = user.Username
		response.ID = user.ID
		response.Email = user.Email
		response.Verified = user.Verified
		c.JSON(http.StatusOK, response)
	}
}
