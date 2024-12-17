package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "request does not contain an access token"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(401, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		//tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		//
		//claims, err := wat.ValidateToken(tokenString)
		//if err != nil {
		//	c.JSON(401, gin.H{"error": err.Error()})
		//	c.Abort()
		//	return
		//}
		//
		//user, err := wat.GetUser(claims.Email)
		//
		//c.Set("user", user)
		c.Next()
	}
}
