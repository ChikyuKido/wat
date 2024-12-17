package wat

type User struct {
	ID          uint         `gorm:"primaryKey"`
	Email       string       `gorm:"size:100;unique;not null"`
	Username    string       `gorm:"size:50;unique;not null"`
	Password    string       `gorm:"size:255;not null"`
	Verified    bool         `gorm:"default:false"`
	Permissions []Permission `gorm:"many2many:user_permissions;"`
}
