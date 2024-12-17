package entity

type Role struct {
	ID          uint         `gorm:"primaryKey"`
	Name        string       `gorm:"size:50;unique;not null"`
	Permissions []Permission `gorm:"many2many:role_permissions"`
}
