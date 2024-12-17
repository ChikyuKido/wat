package wat

import (
	"github.com/ChikyuKido/wat/wat/server/db"
	wat "github.com/ChikyuKido/wat/wat/server/db/entity"
	"github.com/sirupsen/logrus"
)

func InsertNewUser(username, password, email string) bool {
	user := wat.User{
		Email:    email,
		Username: username,
		Password: password,
		Verified: false,
	}
	if err := db.DB().Create(&user).Error; err != nil {
		logrus.Errorf("failed to create use: %v", err)
		return false
	}
	return true
}

func AddPermissionToUser(userID, permissionID uint) bool {
	var user wat.User
	if err := db.DB().First(&user, wat.User{ID: userID}).Error; err != nil {
		logrus.Errorf("failed to get user with id %d: %v", userID, err)
		return false
	}
	var permission wat.Permission
	if err := db.DB().First(&permission, wat.Permission{ID: permissionID}).Error; err != nil {
		logrus.Errorf("failed to get permission with id %d: %v", permissionID, err)
		return false
	}
	if err := db.DB().Model(&user).Association("Permissions").Append(&permission); err != nil {
		logrus.Errorf("failed to append permission with id %d: %v", permissionID, err)
		return false
	}
	return true
}

func AddRoleToUser(userID, roleID uint) bool {
	var user wat.User
	if err := db.DB().First(&user, wat.User{ID: userID}).Error; err != nil {
		logrus.Errorf("failed to get user with id %d: %v", userID, err)
		return false
	}
	var role wat.Role
	if err := db.DB().Preload("Permissions").First(&role, wat.Role{ID: roleID}).Error; err != nil {
		logrus.Errorf("failed to get role with id %d: %v", roleID, err)
		return false
	}
	for _, permission := range role.Permissions {
		if err := db.DB().Model(&user).Association("Permissions").Append(&permission); err != nil {
			logrus.Errorf("failed to append permission with id %d: %v", permission.ID, err)
			return false
		}
	}
	return true
}

func RemoveRoleFromUser(userID, roleID uint) bool {
	var user wat.User
	if err := db.DB().First(&user, wat.User{ID: userID}).Error; err != nil {
		logrus.Errorf("failed to get user with id %d: %v", userID, err)
		return false
	}
	var role wat.Role
	if err := db.DB().Preload("Permissions").First(&role, wat.Role{ID: roleID}).Error; err != nil {
		logrus.Errorf("failed to get role with id %d: %v", roleID, err)
		return false
	}
	for _, permission := range role.Permissions {
		if err := db.DB().Model(&user).Association("Permissions").Delete(&permission); err != nil {
			logrus.Errorf("failed to append permission with id %d: %v", permission.ID, err)
			return false
		}
	}
	return true
}
