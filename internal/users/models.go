package users

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null"`
	Email    string `gorm:"type:varchar(100);unique;not null"`
	Password string `gorm:"not null"`
	Phone    string `gorm:"type:varchar(15);unique;not null"`
	Role     string `gorm:"type:varchar(20);default:'client'"` // Роли: client, owner, admin
}
