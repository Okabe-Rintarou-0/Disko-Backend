package model

import (
	"database/sql"
	"gorm.io/gorm"
)

type Share struct {
	gorm.Model
	UUID     string       `json:"uuid" gorm:"index;type:varchar(36);not null"`
	Password *string      `json:"password" gorm:"type:char(4)"`
	ExpireAt sql.NullTime `json:"expireAt"`
	FileID   uint         `json:"fileID" gorm:"not null"`
	File     File         `json:"file"`
	UserID   uint         `json:"userID" gorm:"not null"`
	User     User         `json:"user"`
}
