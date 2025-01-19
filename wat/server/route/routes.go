package wat

import (
	middleware "github.com/ChikyuKido/wat/wat/server/middleware"
	adminroute "github.com/ChikyuKido/wat/wat/server/route/admin"
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
	user := r.Group("/api/v1/user")
	user.GET("/profile", middleware.RequiredPermission("profile"), userroute.GetProfile())
	admin := r.Group("/api/v1/admin")
	admin.GET("/permissions/list", middleware.RequiredPermission("queryPermissions"), adminroute.GetPermissions())
	admin.GET("/users/list", middleware.RequiredPermission("queryUsers"), adminroute.GetUsers())
	admin.GET("/roles/list", middleware.RequiredPermission("queryRoles"), adminroute.GetRoles())
	admin.POST("/users/addPermissionToUser", middleware.RequiredPermission("changeUserPermissions"), adminroute.AddPermissionToUser())
	admin.POST("/users/addRoleToUser", middleware.RequiredPermission("changeUserPermissions"), adminroute.AddRoleToUser())
	admin.DELETE("/users/deleteUser", middleware.RequiredPermission("deleteUser"), adminroute.DeleteUser())
	admin.DELETE("/users/removePermissionFromUser", middleware.RequiredPermission("changeUserPermissions"), adminroute.RemovePermissionToUser())
	admin.DELETE("/users/removeRoleFromUser", middleware.RequiredPermission("changeUserPermissions"), adminroute.RemoveRoleFromUser())
}
