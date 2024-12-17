package wat

import (
	entity "github.com/ChikyuKido/wat/wat/server/db/entity"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// GetUserFromContext returns the user from the context. If the context does not contain a user it prints an error and adds an error message to the response.
// please note that this shouldn't happen in any case
func GetUserFromContext(c *gin.Context) *entity.User {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not validate user. Please try again later"})
		c.Abort()
		logrus.Errorf("Could not retrieve user from context. This is a serious problem")
		return nil
	}
	u, ok := user.(*entity.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not validate user. Please try again later"})
		c.Abort()
		logrus.Errorf("Could not retrieve user from context. This is a serious problem")
		return nil
	}
	return u
}
