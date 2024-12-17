package repo

import (
	"Quotium/internal/server/db"
	"Quotium/internal/server/db/entity"
	"github.com/sirupsen/logrus"
)

func InsertNewPermission(name string) bool {
	permission := entity.Permission{
		Name: name,
	}
	if err := db.DB().Create(&permission).Error; err != nil {
		logrus.Errorf("failed to insert new permission: %v", err)
		return false
	}
	return true
}
