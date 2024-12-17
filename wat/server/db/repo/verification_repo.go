package wat

import (
	db "github.com/ChikyuKido/wat/wat/server/db"
	entity "github.com/ChikyuKido/wat/wat/server/db/entity"
	"github.com/sirupsen/logrus"
	"time"
)

func InsertNewVerification(uuid string, userID uint) bool {
	verification := entity.Verification{
		UUID:    uuid,
		UserID:  userID,
		Expires: time.Now().Unix() + 60*15,
	}
	if err := db.DB().Create(&verification).Error; err != nil {
		logrus.Errorf("Failed to create a verification: %v", err)
		return false
	}
	return true
}
func GetVerificationFromUUID(uuid string) *entity.Verification {
	var verification entity.Verification
	if err := db.DB().Where(entity.Verification{UUID: uuid}).First(&verification).Error; err != nil {
		return nil
	}
	return &verification
}
func DeleteVerificationByUUID(uuid string) {
	if err := db.DB().Delete(&entity.Verification{}, uuid).Error; err != nil {
		logrus.Errorf("Failed to delete a verification: %v", err)
	}
}

func CountUserVerifications(userID uint) int {
	var verifications []entity.Verification
	if err := db.DB().Where(entity.Verification{UserID: userID}).Find(&verifications).Error; err != nil {
		logrus.Errorf("Failed to count verifications: %v", err)
		return 0
	}
	return len(verifications)
}
