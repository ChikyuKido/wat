package admin

import (
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPermissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions := repo.GetAllPermissions()
		c.JSON(http.StatusOK, permissions)
	}
}
