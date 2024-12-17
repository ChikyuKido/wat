package entity

type Verification struct {
	UUID    string `gorm:"primaryKey"`
	UserID  uint
	User    User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Expires int64 `gorm:"type:datetime;not null"`
}
