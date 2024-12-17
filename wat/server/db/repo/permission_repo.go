package wat

import (
	db "github.com/ChikyuKido/wat/wat/server/db"
	wat "github.com/ChikyuKido/wat/wat/server/db/entity"
	"github.com/sirupsen/logrus"
)

func InsertNewPermission(name string) bool {
	permission := wat.Permission{
		Name: name,
	}
	if err := db.DB().Create(&permission).Error; err != nil {
		logrus.Errorf("failed to insert new permission: %v", err)
		return false
	}
	return true
}
