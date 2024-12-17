package wat

type Permission struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:50;unique;not null"`
}
