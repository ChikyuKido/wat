package wat

import (
	db "github.com/ChikyuKido/wat/wat/server/db"
	wat "github.com/ChikyuKido/wat/wat/server/db/entity"
	"github.com/sirupsen/logrus"
	"time"
)

func InsertNewVerification(uuid string, userID uint) bool {
	verification := wat.Verification{
		UUID:    uuid,
		UserID:  userID,
		Expires: time.Now().Unix() + 3600,
	}
	if err := db.DB().Create(&verification).Error; err != nil {
		logrus.Errorf("Failed to create a verification: %v", err)
		return false
	}
	return true
}
