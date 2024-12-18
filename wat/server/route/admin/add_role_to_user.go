package admin

import (
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddRoleToUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData = struct {
			UserID                  uint `json:"user_id"`
			RoleID                  uint `json:"role_id"`
			OverwriteOldPermissions bool `json:"overwrite_old_permissions"`
		}{}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad json format"})
			return
		}
		if requestData.UserID == 0 || requestData.RoleID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad json format"})
			return
		}
		if !repo.DoesUserByIDExists(requestData.UserID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "user does not exist"})
			return
		}
		if !repo.DoesRoleByIDExists(requestData.RoleID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "role does not exist"})
			return
		}
		if requestData.OverwriteOldPermissions {
			if !repo.RemoveAllPermissionsFromUser(requestData.UserID) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove old permissions"})
				return
			}
		}

		if !repo.AddRoleToUser(requestData.UserID, requestData.RoleID) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign role to user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "added role to user"})
	}
}
