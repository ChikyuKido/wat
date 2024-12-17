package wat

import (
	"github.com/ChikyuKido/wat/wat/server/db/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabase(passedDB *gorm.DB) {
	db = passedDB
	err := db.AutoMigrate(wat.Permission{}, wat.Role{}, wat.User{}, wat.Verification{})
	if err != nil {
		logrus.Fatalf("failed to migrate database: %v", err)
	}
	logrus.Info("Database initialized")
}

func DB() *gorm.DB {
	return db
}
