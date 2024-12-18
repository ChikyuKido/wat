package admin

import (
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRoles() gin.HandlerFunc {
	return func(c *gin.Context) {
		roles := repo.GetAllRoles()
		c.JSON(http.StatusOK, roles)
	}
}
