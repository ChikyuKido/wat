package middleware

import (
	wat "github.com/ChikyuKido/wat/wat/server/db/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RequiredPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		u, ok := user.(*wat.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
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
