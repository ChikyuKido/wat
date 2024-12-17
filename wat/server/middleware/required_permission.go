package wat

import (
	wat "github.com/ChikyuKido/wat/wat/server/db/entity"
	util "github.com/ChikyuKido/wat/wat/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RequiredPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := util.GetUserFromContext(c)
		if u == nil {
			c.Abort()
			return
		}
		if hasPermission(u, permission) {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"message": "You do not have permission to access this endpoint"})
			c.Abort()
		}
	}
}

func hasPermission(user *wat.User, roleToCheck string) bool {
	for _, permission := range user.Permissions {
		if strings.TrimSpace(permission.Name) == roleToCheck {
			return true
		}
	}
	return false
}
