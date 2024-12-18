package wat

import (
	"errors"
	db "github.com/ChikyuKido/wat/wat/server/db"
	entity "github.com/ChikyuKido/wat/wat/server/db/entity"
	util "github.com/ChikyuKido/wat/wat/util"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InsertNewUser(username, password, email string) bool {
	user := entity.User{
		Email:    email,
		Username: username,
		Password: password,
		Verified: !util.Config.EmailVerification,
	}
	if err := db.DB().Create(&user).Error; err != nil {
		logrus.Errorf("failed to create use: %v", err)
		return false
	}
	return true
}
func VerifyUser(userid uint) bool {
	var user entity.User
	if err := db.DB().First(&user, entity.User{ID: userid}).Error; err != nil {
		return false
	}
	user.Verified = true
	if err := db.DB().Save(&user).Error; err != nil {
		return false
	}
	return true
}
func GetAllUsers() []entity.User {
	var users []entity.User
	if err := db.DB().Preload("Permissions").Find(&users).Error; err != nil {
		logrus.Errorf("failed to get all users: %v", err)
		return nil
	}
	return users
}
func GetUserByEmail(email string) *entity.User {
	var user entity.User
	if err := db.DB().Preload("Permissions").First(&user, entity.User{Email: email}).Error; err != nil {
		logrus.Errorf("failed to get user: %v", err)
		return nil
	}
	return &user
}
func GetUserByUsername(username string) *entity.User {
	var user entity.User
	if err := db.DB().Preload("Permissions").First(&user, entity.User{Username: username}).Error; err != nil {
		logrus.Errorf("failed to get user: %v", err)
		return nil
	}
	return &user
}
func GetUserByID(id uint) *entity.User {
	var user entity.User
	if err := db.DB().Preload("Permissions").Where(entity.User{ID: id}).First(&user).Error; err != nil {
		logrus.Errorf("failed to get user: %v", err)
		return nil
	}
	return &user
}
func DoesUserByIDExists(id uint) bool {
	var user entity.User
	if err := db.DB().First(&user, entity.User{ID: id}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		logrus.Errorf("failed to get user: %v", err)
		return true
	}
	return true
}
func DoesUserByEmailExist(email string) bool {
	var user entity.User
	err := db.DB().First(&user, entity.User{Email: email}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		logrus.Errorf("failed to get user: %v", err)
		return true
	}
	return true
}
func DoesUserByUsernameExist(username string) bool {
	var user entity.User
	err := db.DB().First(&user, entity.User{Username: username}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		logrus.Errorf("failed to get user: %v", err)
		return true
	}
	return true
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

func RemoveAllPermissionsFromUser(userID uint) bool {
	var user entity.User
	if err := db.DB().Preload("Permissions").First(&user, entity.User{ID: userID}).Error; err != nil {
		logrus.Errorf("failed to get user with id %d: %v", userID, err)
		return false
	}
	for _, permission := range user.Permissions {
		if err := db.DB().Model(&user).Association("Permissions").Delete(&permission); err != nil {
			logrus.Errorf("failed to append permission with id %d: %v", permission.ID, err)
			return false
		}
	}
	return true
}

func DeleteUser(userID uint) bool {
	var user = entity.User{
		ID: userID,
	}
	if err := db.DB().Delete(&user).Error; err != nil {
		logrus.Errorf("failed to delete user: %v", err)
		return false
	}
	return true
}
func RemovePermissionFromUser(userID, permissionID uint) bool {
	var user entity.User
	if err := db.DB().First(&user, entity.User{ID: userID}).Error; err != nil {
		logrus.Errorf("failed to get user with id %d: %v", userID, err)
		return false
	}
	if err := db.DB().Model(&user).Association("Permissions").Delete(entity.Permission{ID: permissionID}); err != nil {
		logrus.Errorf("failed to remove permission with id %d: %v", permissionID, err)
		return false
	}
	return true
}
