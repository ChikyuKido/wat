package wat

import (
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Verify() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("uuid")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
			return
		}
		verification := repo.GetVerificationFromUUID(id)
		if verification == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "uuid not found"})
			return
		}
		if verification.Expires < time.Now().Unix() {
			c.JSON(http.StatusForbidden, gin.H{"error": "verification expired"})
			return
		}
		if !repo.VerifyUser(verification.UserID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid user"})
			return
		}
		repo.RemoveAllPermissionsFromUser(verification.UserID)
		repo.AddRoleToUser(verification.UserID, repo.GetRoleByName("user").ID) // user role
		repo.DeleteVerificationByUUID(id)
		c.JSON(http.StatusOK, gin.H{"message": "successful verified user"})
	}
}
