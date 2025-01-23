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
			if firstRegister(c) {
				return
			}
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

		if firstRegister(c) {
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func firstRegister(c *gin.Context) bool {
	if util.Config.FirstUser {
		if c.FullPath() == "/auth/adminRegister" || c.FullPath() == "/api/v1/auth/register" {
			c.Next()
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "/auth/adminRegister")
			c.Abort()
		}
		return true
	} else {
		if c.FullPath() == "/auth/adminRegister" {
			c.Redirect(http.StatusTemporaryRedirect, "/notfound")
			c.Abort()
			return true
		}
	}
	return false
}
