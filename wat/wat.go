package wat

import (
	db "github.com/ChikyuKido/wat/wat/server/db"
	repo "github.com/ChikyuKido/wat/wat/server/db/repo"
	middleware "github.com/ChikyuKido/wat/wat/server/middleware"
	wat "github.com/ChikyuKido/wat/wat/server/route"
	static "github.com/ChikyuKido/wat/wat/server/static"
	util "github.com/ChikyuKido/wat/wat/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"strings"
)

var Roles = map[string][]string{
	"guest":          {"login", "register"},
	"unverifiedUser": {"login", "register", "sendVerification", "profile"},
	"user":           {"login", "register", "sendVerification", "profile"},
	"admin":          {"login", "register", "sendVerification", "profile", "dashboard", "queryPermissions", "queryUsers", "queryRoles", "changeUserPermissions", "deleteUser"},
}

func InitWat(engine *gin.Engine, database *gorm.DB, firstInit bool) {
	initEnv()
	logrus.Info("checked env variables")
	db.InitDatabase(database)
	logrus.Info("initialized database")
	if firstInit {
		dataInit()
		logrus.Info("first init. Adding default data")
	}
	util.Config.FirstUser = len(repo.GetAllUsers()) == 1
	wat.InitRoutes(engine)
	logrus.Info("initialized routes")
}

func InitWatWebsite(engine *gin.Engine, basePath string) {
	sites := engine.Group("/")
	sitesWithoutAuth := engine.Group("/")
	sites.Use(middleware.AuthMiddleware())
	static.ServeFolder("/css/", basePath+"/css", nil, "wat", sitesWithoutAuth, "")
	static.ServeFolder("/js/", basePath+"/js", nil, "wat", sitesWithoutAuth, "")
	sites.GET("/admin/dashboard", middleware.RequiredPermission("dashboard", true), static.ServeFile(basePath+"/html/admin/dashboard.html", nil, "wat"))
	sites.GET("/auth/login", middleware.RequiredPermission("login", true), static.ServeFile(basePath+"/html/auth/login.html", nil, "wat"))
	sites.GET("/auth/register", middleware.RequiredPermission("register", true), static.ServeFile(basePath+"/html/auth/register.html", nil, "wat"))
	sites.GET("/auth/adminRegister", static.ServeFile(basePath+"/html/auth/adminRegister.html", nil, "wat"))
	sites.GET("/auth/verify", static.ServeFile(basePath+"/html/auth/verify.html", nil, "wat"))
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
	util.Config.AllowedEmails = strings.Split(os.Getenv("ALLOWED_EMAILS"), ",")
	debug := os.Getenv("DEBUG")
	checkEnv(debug, "DEBUG")
	if debug == "true" {
		util.Config.Debug = true
	} else {
		util.Config.Debug = false
	}

	util.Config.JwtSecret = os.Getenv("JWT_SECRET")
	checkEnv(util.Config.JwtSecret, "JWT_SECRET")
	util.Config.ResourceVersion = os.Getenv("RESOURCE_VERSION")
	checkEnv(util.Config.ResourceVersion, "RESOURCE_VERSION")
	return true
}
func checkEnv(value string, envName string) {
	if value == "" {
		logrus.Fatalf("Please enter a %s in the envs", envName)
	}
}
func dataInit() {
	var allPermission = make(map[string]uint)
	var counter uint = 1
	for _, permissions := range Roles {
		for _, permission := range permissions {
			if _, exists := allPermission[permission]; exists {
				continue
			}
			allPermission[permission] = counter
			repo.InsertNewPermission(permission)
			counter++
		}
	}
	counter = 1
	for role, permissions := range Roles {
		repo.InsertNewRole(role)
		for _, permission := range permissions {
			repo.AddPermissionToRole(counter, allPermission[permission])
		}
		counter++
	}
	repo.InsertNewUser("guest", "", "guest")
	repo.AddRoleToUser(1, repo.GetRoleByName("guest").ID)
}
