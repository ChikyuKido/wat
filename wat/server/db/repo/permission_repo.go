package wat

import (
	"errors"
	db "github.com/ChikyuKido/wat/wat/server/db"
	entity "github.com/ChikyuKido/wat/wat/server/db/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

func GetPermissionByID(id uint) *entity.Permission {
	var permission entity.Permission
	if err := db.DB().Where(entity.Permission{ID: id}).First(&permission).Error; err != nil {
		logrus.Errorf("failed to get permission: %v", err)
		return nil
	}
	return &permission
}
func GetPermissionByName(name string) *entity.Permission {
	var permission entity.Permission
	if err := db.DB().Where("name = ?", name).First(&permission).Error; err != nil {
		logrus.Errorf("failed to get permission: %v", err)
		return nil
	}
	return &permission
}
func DoesPermissionByIDExists(id uint) bool {
	var permission entity.Permission
	if err := db.DB().Where(entity.Permission{ID: id}).First(&permission).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		logrus.Errorf("failed to get permission: %v", err)
		return true
	}
	return true
}

func GetAllPermissions() []entity.Permission {
	var permissions []entity.Permission
	if err := db.DB().Find(&permissions).Error; err != nil {
		logrus.Errorf("failed to query all permissions: %v", err)
		return nil
	}
	return permissions
}
