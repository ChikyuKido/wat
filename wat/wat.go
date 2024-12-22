package wat

import (
	db "github.com/ChikyuKido/wat/wat/server/db"
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	middleware "github.com/ChikyuKido/wat/wat/server/middleware"
	wat "github.com/ChikyuKido/wat/wat/server/route"
	util "github.com/ChikyuKido/wat/wat/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
)

func InitWat(engine *gin.Engine, database *gorm.DB, firstInit bool) {
	initEnv()
	logrus.Info("checked env variables")
	db.InitDatabase(database)
	logrus.Info("initialized database")
	if firstInit {
		dataInit()
		logrus.Info("first init. Adding default data")
	}
	engine.Use(middleware.AuthMiddleware())
	wat.InitRoutes(engine)
	logrus.Info("initialized routes")
}

func InitWatWebsite(engine *gin.Engine, basePath string) {
	sites := engine.Group("/")
	util.ServeFolder("/css/", basePath+"/css", sites)
	util.ServeFolder("/js/", basePath+"/js", sites)
	util.ServeFile("/admin/dashboard", basePath+"/html/admin/dashboard.html", sites)
	util.ServeFile("/auth/login", basePath+"/html/auth/login.html", sites)
	util.ServeFile("/auth/register", basePath+"/html/auth/register.html", sites)
}

func initEnv() bool {
	util.Config.WatBaseURL = os.Getenv("WAT_BASEURL")
	checkEnv(util.Config.WatBaseURL, "WAT_BASEURL")

	emailVerification := os.Getenv("EMAIL_VERIFICATION")
	checkEnv(emailVerification, "EMAIL_VERIFICATION")
	if emailVerification == "true" {
		util.Config.EmailVerification = true
		util.Config.SmtpServer = os.Getenv("SMTP_SERVER")
		checkEnv(util.Config.SmtpServer, "SMTP_SERVER")
		util.Config.SmtpHost = os.Getenv("SMTP_HOST")
		checkEnv(util.Config.SmtpHost, "SMTP_HOST")
		util.Config.SmtpPassword = os.Getenv("SMTP_PASSWORD")
		checkEnv(util.Config.SmtpPassword, "SMTP_PASSWORD")
		util.Config.SmtpEmail = os.Getenv("SMTP_EMAIL")
		checkEnv(util.Config.SmtpEmail, "SMTP_EMAIL")
		util.Config.EmailVerificationUrl = os.Getenv("EMAIL_VERIFICATION_URL")
		checkEnv(util.Config.EmailVerificationUrl, "EMAIL_VERIFICATION_URL")
	} else {
		util.Config.EmailVerification = false
	}
	debug := os.Getenv("DEBUG")
	checkEnv(debug, "DEBUG")
	if debug == "true" {
		util.Config.Debug = true
	} else {
		util.Config.Debug = false
	}

	util.Config.JwtSecret = os.Getenv("JWT_SECRET")
	checkEnv(util.Config.JwtSecret, "JWT_SECRET")

	return true
}
func checkEnv(value string, envName string) {
	if value == "" {
		logrus.Fatalf("Please enter a %s in the envs", envName)
	}
}
func dataInit() {
	repo.InsertNewPermission("login")
	repo.InsertNewPermission("register")
	repo.InsertNewPermission("sendVerification")
	repo.InsertNewPermission("queryPermissions")
	repo.InsertNewPermission("queryUsers")
	repo.InsertNewPermission("queryRoles")
	repo.InsertNewPermission("changeUserPermissions")
	repo.InsertNewPermission("deleteUser")
	repo.InsertNewRole("guest")
	repo.AddPermissionToRole(1, 1)
	repo.AddPermissionToRole(1, 2)
	repo.InsertNewRole("unverifiedUser")
	repo.AddPermissionToRole(2, 1)
	repo.AddPermissionToRole(2, 2)
	repo.AddPermissionToRole(2, 3)
	repo.InsertNewRole("user")
	repo.AddPermissionToRole(3, 1)
	repo.AddPermissionToRole(3, 2)
	repo.AddPermissionToRole(3, 3)
	repo.InsertNewUser("guest", "", "guest")
	repo.AddRoleToUser(1, 1)
}
