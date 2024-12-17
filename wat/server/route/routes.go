package wat

import (
	middleware "github.com/ChikyuKido/wat/wat/server/middleware"
	userroute "github.com/ChikyuKido/wat/wat/server/route/user"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	auth := r.Group("/api/v1/auth")
	auth.POST("/login", userroute.Login(), middleware.RequiredPermission("login"))
	auth.POST("/register", userroute.Login(), middleware.RequiredPermission("register"))
}
