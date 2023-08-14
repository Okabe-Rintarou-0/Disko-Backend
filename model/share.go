package model

import (
	"gorm.io/gorm"
	"time"
)

type Share struct {
	gorm.Model
	UUID     string    `json:"uuid" gorm:"index;type:varchar(36);not null"`
	Password string    `json:"password" gorm:"type:char(4)"`
	ExpireAt time.Time `json:"expireAt" gorm:"not null"`
	FileID   uint      `json:"fileID" gorm:"not null"`
	File     File      `json:"file"`
	UserID   uint      `json:"userID" gorm:"not null"`
	User     User      `json:"user"`
}
