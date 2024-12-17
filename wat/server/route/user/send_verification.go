package wat

import (
	"fmt"
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	util "github.com/ChikyuKido/wat/wat/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		id := uuid.New()
		if !repo.InsertNewVerification(id.String(), user.ID) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred. Please try again later"})
			return
		}
		err := util.SendEmail("Email verification", fmt.Sprintf(`
				Please click on the link below to verify your Email.
				http://localhost:8080/api/auth/verify?uuid=%s
				This link will expire in 15 minutes
		`, id), user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "sent verification link to email"})

	}
}
