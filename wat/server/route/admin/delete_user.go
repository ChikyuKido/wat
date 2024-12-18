package admin

import (
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData = struct {
			UserID uint `json:"user_id"`
		}{}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad json format"})
			return
		}
		if requestData.UserID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad json format"})
			return
		}
		if !repo.DeleteUser(requestData.UserID) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "successfully deleted user"})
	}
}
