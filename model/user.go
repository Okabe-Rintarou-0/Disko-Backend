package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string `json:"name" gorm:"index;type:varchar(100);collate:utf8;not null"`
	// encrypted
	Password string `json:"password" gorm:"type:char(60);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);not null"`
	Files    []File `json:"files" gorm:"foreignKey:owner"`
}
