package admin

import (
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RemovePermissionToUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData = struct {
			UserID       uint `json:"user_id"`
			PermissionID uint `json:"permission_id"`
		}{}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad json format"})
			return
		}
		if requestData.UserID == 0 || requestData.PermissionID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad json format"})
			return
		}

		if !repo.DoesPermissionByIDExists(requestData.PermissionID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission does not exist"})
			return
		}
		if !repo.DoesUserByIDExists(requestData.UserID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "user does not exist"})
			return
		}
		if !repo.RemovePermissionFromUser(requestData.UserID, requestData.PermissionID) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove permission from user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "removed permission from user"})
	}
}
