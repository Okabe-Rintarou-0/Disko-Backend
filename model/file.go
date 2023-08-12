package model

import (
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Name     string `json:"name" gorm:"index;type:varchar(100);not null"`
	Ext      string `json:"ext" gorm:"index;type:char(10);not null"`
	Size     int64  `json:"size" gorm:"not null"`
	UUID     string `json:"uuid" gorm:"type:varchar(36);not null"`
	Path     string `json:"path" gorm:"type:varchar(255);not null"`
	Owner    uint   `json:"owner" gorm:"index;not null"`
	IsDir    bool   `json:"isDir" gorm:"not null"`
	Private  bool   `json:"private" gorm:"not null"`
	ParentID *uint
	Parent   *File `gorm:"foreignkey:ParentID"`
}
