package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string `json:"name" gorm:"index;type:varchar(100);not null"`
	// encrypted
	Password string `json:"password" gorm:"type:varchar(32);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);not null"`
}
