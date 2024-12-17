package repo

import (
	"Quotium/internal/server/db"
	"Quotium/internal/server/db/entity"
	"github.com/sirupsen/logrus"
)

func InsertNewRole(name string) bool {
	role := entity.Role{Name: name}
	if err := db.DB().Create(&role).Error; err != nil {
		logrus.Errorf("failed to insert new permission: %v", err)
		return false
	}
	return true
}

func AddPermissionToRole(roleId uint, permissionID uint) bool {
	var role entity.Role
	if err := db.DB().First(&role, entity.Role{ID: roleId}).Error; err != nil {
		logrus.Errorf("failed to get role with id %d: %v", roleId, err)
		return false
	}
	var permission entity.Permission
	if err := db.DB().First(&permission, entity.Permission{ID: permissionID}).Error; err != nil {
		logrus.Errorf("failed to get permission with id %d: %v", permissionID, err)
		return false
	}
	if err := db.DB().Model(&role).Association("Permissions").Append(&permission); err != nil {
		logrus.Errorf("failed to append permission with id %d: %v", permissionID, err)
		return false
	}
	return true
}
