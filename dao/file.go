package dao

import (
	"disko/model"
	"disko/repository/query"
)

type IFileDAO interface {
	FindById(id uint) (*model.File, error)
	FindByUUID(uuid string) (*model.File, error)
	Save(files ...*model.File) error
}

func NewFileDAO() IFileDAO {
	return &FileDAO{}
}

type FileDAO struct{}

func (f *FileDAO) FindById(id uint) (*model.File, error) {
	u := query.Use(db).File
	return u.WithContext(ctx).Where(u.ID.Eq(id)).Take()
}

func (f *FileDAO) FindByUUID(uuid string) (*model.File, error) {
	u := query.Use(db).File
	return u.WithContext(ctx).Where(u.UUID.Eq(uuid)).Take()
}

func (f *FileDAO) Save(files ...*model.File) error {
	u := query.Use(db).File
	return u.WithContext(ctx).Save(files...)
}
