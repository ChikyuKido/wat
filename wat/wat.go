package wat

import (
	db "github.com/ChikyuKido/wat/wat/server/db"
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	middleware "github.com/ChikyuKido/wat/wat/server/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitWat(engine *gin.Engine, database *gorm.DB, firstInit bool) {
	db.InitDatabase(database)

	if firstInit {
		repo.InsertNewPermission("login")
		repo.InsertNewPermission("register")
		repo.InsertNewPermission("verify")
		repo.InsertNewPermission("sendVerification")
		repo.InsertNewRole("guest")
		repo.AddPermissionToRole(1, 1)
		repo.AddPermissionToRole(1, 2)
		repo.InsertNewRole("unverifiedUser")
		repo.AddPermissionToRole(2, 1)
		repo.AddPermissionToRole(2, 2)
		repo.AddPermissionToRole(2, 3)
		repo.AddPermissionToRole(2, 4)
		repo.InsertNewUser("guest", "", "guest")
		repo.AddRoleToUser(1, 1)
	}
	engine.Use(middleware.AuthMiddleware())
}
