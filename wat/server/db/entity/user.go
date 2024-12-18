package wat

type User struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Email       string       `gorm:"size:100;unique;not null" json:"email"`
	Username    string       `gorm:"size:50;unique;not null" json:"username"`
	Password    string       `gorm:"size:255;not null" json:"password"`
	Verified    bool         `gorm:"default:false" json:"verified"`
	Permissions []Permission `gorm:"many2many:user_permissions;" json:"permissions"`
}
