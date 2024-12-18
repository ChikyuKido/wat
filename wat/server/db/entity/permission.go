package wat

type Permission struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:50;unique;not null" json:"name"`
}
