package wat

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerData = struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			Password string `json:"password"`
		}{}
		if err := c.ShouldBind(&registerData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request data"})
			return
		}

	}

}
