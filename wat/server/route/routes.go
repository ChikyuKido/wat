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
	auth.POST("/login", middleware.RequiredPermission("login", false), userroute.Login())
	auth.POST("/register", middleware.RequiredPermission("register", false), userroute.Register())
	if util.Config.EmailVerification {
		auth.POST("/sendVerification", middleware.RequiredPermission("sendVerification", false), userroute.SendVerification())
		auth.POST("/verify", userroute.Verify())
	}
	user := r.Group("/api/v1/user")
	user.GET("/profile", middleware.RequiredPermission("profile", false), userroute.GetProfile())
	admin := r.Group("/api/v1/admin")
	admin.GET("/permissions/list", middleware.RequiredPermission("queryPermissions", false), adminroute.GetPermissions())
	admin.GET("/users/list", middleware.RequiredPermission("queryUsers", false), adminroute.GetUsers())
	admin.GET("/roles/list", middleware.RequiredPermission("queryRoles", false), adminroute.GetRoles())
	admin.POST("/users/addPermissionToUser", middleware.RequiredPermission("changeUserPermissions", false), adminroute.AddPermissionToUser())
	admin.POST("/users/addRoleToUser", middleware.RequiredPermission("changeUserPermissions", false), adminroute.AddRoleToUser())
	admin.DELETE("/users/deleteUser", middleware.RequiredPermission("deleteUser", false), adminroute.DeleteUser())
	admin.DELETE("/users/removePermissionFromUser", middleware.RequiredPermission("changeUserPermissions", false), adminroute.RemovePermissionToUser())
	admin.DELETE("/users/removeRoleFromUser", middleware.RequiredPermission("changeUserPermissions", false), adminroute.RemoveRoleFromUser())
}
