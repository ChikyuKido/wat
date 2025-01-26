package wat

import (
	"github.com/ChikyuKido/wat/wat/helper"
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	"github.com/gin-gonic/gin"
	"net/http"
)

var verificationUUIDS = make(map[string]uint)

func SendVerification() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("verificationUUID")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "verificationUUID is required"})
			return
		}
		if _, exists := verificationUUIDS[id]; !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "verificationUUID not found"})
			return
		}
		user := repo.GetUserByID(verificationUUIDS[id])
		if user == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		verificationCount := repo.CountUserVerifications(user.ID)
		if verificationCount > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Too many pending verifications for the user"})
			return
		}
		if !helper.SendEmailVerificationForUser(user) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email verification failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "sent verification link to email"})
		delete(verificationUUIDS, id)
	}
}
