package wat

import (
	"errors"
	db "github.com/ChikyuKido/wat/wat/server/db"
	entity "github.com/ChikyuKido/wat/wat/server/db/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InsertNewRole(name string) bool {
	role := entity.Role{Name: name}
	if err := db.DB().Create(&role).Error; err != nil {
		logrus.Errorf("failed to insert new permission: %v", err)
		return false
	}
	return true
}
func GetAllRoles() []entity.Role {
	var roles []entity.Role
	if err := db.DB().Preload("Permissions").Find(&roles).Error; err != nil {
		logrus.Errorf("failed to query all roles: %v", err)
		return nil
	}
	return roles
}
func GetRoleByName(name string) *entity.Role {
	var role entity.Role
	if err := db.DB().Where(entity.Role{Name: name}).First(&role).Error; err != nil {
		return nil
	}
	return &role
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

func DoesRoleByIDExists(id uint) bool {
	var permission entity.Role
	if err := db.DB().Where(entity.Role{ID: id}).First(&permission).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		logrus.Errorf("failed to get role: %v", err)
		return true
	}
	return true
}
