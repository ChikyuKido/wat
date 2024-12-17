package wat

import (
	middleware "github.com/ChikyuKido/wat/wat/server/middleware"
	userroute "github.com/ChikyuKido/wat/wat/server/route/user"
	util "github.com/ChikyuKido/wat/wat/util"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	auth := r.Group("/api/v1/auth")
	auth.POST("/login", middleware.RequiredPermission("login"), userroute.Login())
	auth.POST("/register", middleware.RequiredPermission("register"), userroute.Register())
	if util.Config.EmailVerification {
		auth.POST("/sendVerification", middleware.RequiredPermission("sendVerification"), userroute.SendVerification())
		auth.POST("/verify", userroute.Verify())
	}
}
