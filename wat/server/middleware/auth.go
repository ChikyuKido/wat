package wat

import (
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	util "github.com/ChikyuKido/wat/wat/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, _ := c.Cookie("jwt")
		// guest login
		if tokenString == "" || tokenString == "guest" {
			guest := repo.GetUserByEmail("guest")
			if guest == nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No guest user found."})
				c.Abort()
				return
			}
			c.Set("user", guest)
			c.Next()
			return
		}

		token, err := util.GetToken(tokenString)
		if err != nil || !token.Valid {
			c.SetCookie("jwt", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.SetCookie("jwt", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		email := claims["email"].(string)
		user := repo.GetUserByEmail(email)
		if user == nil {
			c.SetCookie("jwt", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
