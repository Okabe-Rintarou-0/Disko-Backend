package dao

import (
	"disko/model"
	"disko/repository/query"
)

type IShareDAO interface {
	FindByID(id uint) (*model.Share, error)
	FindByUUIDWithFile(uuid string) (*model.Share, error)
	FindByUUIDWithFileAndUser(uuid string) (*model.Share, error)
	Save(shares ...*model.Share) error
}

func NewShareDAO() IShareDAO {
	return &ShareDAO{}
}

type ShareDAO struct{}

func (sd *ShareDAO) FindByUUIDWithFileAndUser(uuid string) (*model.Share, error) {
	s := query.Use(db).Share
	return s.WithContext(ctx).Preload(s.File).Preload(s.User).Where(s.UUID.Eq(uuid)).Take()
}

func (sd *ShareDAO) FindByUUIDWithFile(uuid string) (*model.Share, error) {
	s := query.Use(db).Share
	return s.WithContext(ctx).Preload(s.File).Where(s.UUID.Eq(uuid)).Take()
}

func (sd *ShareDAO) FindByID(id uint) (*model.Share, error) {
	s := query.Use(db).Share
	return s.WithContext(ctx).Where(s.ID.Eq(id)).Take()
}

func (sd *ShareDAO) Save(shares ...*model.Share) error {
	s := query.Use(db).Share
	return s.WithContext(ctx).Save(shares...)
}
