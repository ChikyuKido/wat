package wat

type Verification struct {
	UUID    string `gorm:"primaryKey" json:"uuid"`
	UserID  uint
	User    User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	Expires int64 `gorm:"not null" json:"expires"`
}
