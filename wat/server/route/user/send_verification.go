package wat

import (
	"github.com/ChikyuKido/wat/wat/helper"
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	util "github.com/ChikyuKido/wat/wat/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendVerification() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := util.GetUserFromContext(c)
		if user == nil {
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

	}
}
