package wat

import (
	db "github.com/ChikyuKido/wat/wat/server/db"
	entity "github.com/ChikyuKido/wat/wat/server/db/entity"
	"github.com/sirupsen/logrus"
)

func InsertNewUser(username, password, email string) bool {
	user := entity.User{
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

func GetUserByEmail(email string) *entity.User {
	var user entity.User
	if err := db.DB().First(&user, entity.User{Email: email}).Error; err != nil {
		logrus.Errorf("failed to get user: %v", err)
		return nil
	}
	return &user
}
func GetUserByUsername(username string) *entity.User {
	var user entity.User
	if err := db.DB().First(&user, entity.User{Username: username}).Error; err != nil {
		logrus.Errorf("failed to get user: %v", err)
		return nil
	}
	return &user
}
func DoesUserByEmailExist(email string) bool {
	return GetUserByEmail(email) != nil
}
func DoesUserByUsernameExist(email string) bool {
	return GetUserByUsername(email) != nil
}

func AddPermissionToUser(userID, permissionID uint) bool {
	var user entity.User
	if err := db.DB().First(&user, entity.User{ID: userID}).Error; err != nil {
		logrus.Errorf("failed to get user with id %d: %v", userID, err)
		return false
	}
	var permission entity.Permission
	if err := db.DB().First(&permission, entity.Permission{ID: permissionID}).Error; err != nil {
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
	var user entity.User
	if err := db.DB().First(&user, entity.User{ID: userID}).Error; err != nil {
		logrus.Errorf("failed to get user with id %d: %v", userID, err)
		return false
	}
	var role entity.Role
	if err := db.DB().Preload("Permissions").First(&role, entity.Role{ID: roleID}).Error; err != nil {
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
	var user entity.User
	if err := db.DB().First(&user, entity.User{ID: userID}).Error; err != nil {
		logrus.Errorf("failed to get user with id %d: %v", userID, err)
		return false
	}
	var role entity.Role
	if err := db.DB().Preload("Permissions").First(&role, entity.Role{ID: roleID}).Error; err != nil {
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
