package repo

import (
	"Quotium/internal/server/db"
	"Quotium/internal/server/db/entity"
	"github.com/sirupsen/logrus"
	"time"
)

func InsertNewVerification(uuid string, userID uint) bool {
	verification := entity.Verification{
		UUID:    uuid,
		UserID:  userID,
		User:    entity.User{},
		Expires: time.Now().Unix() + 3600,
	}
	if err := db.DB().Create(&verification).Error; err != nil {
		logrus.Errorf("Failed to create a verification: %v", err)
		return false
	}
	return true
}
