package admin

import (
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users := repo.GetAllUsers()
		c.JSON(http.StatusOK, users)
	}
}
