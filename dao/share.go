package dao

import (
	"disko/model"
	"disko/repository/query"
)

type IShareDAO interface {
	FindByID(id uint) (*model.Share, error)
	FindByUUID(uuid string) (*model.Share, error)
	FindByUUIDWithFile(uuid string) (*model.Share, error)
	FindByUUIDWithFileAndUser(uuid string) (*model.Share, error)
	FindByUserId(userId uint) ([]*model.Share, error)
	Save(shares ...*model.Share) error
	DeleteByUUID(uuid string) error
}

func NewShareDAO() IShareDAO {
	return &ShareDAO{}
}

type ShareDAO struct{}

func (sd *ShareDAO) FindByUUID(uuid string) (*model.Share, error) {
	s := query.Use(db).Share
	return s.WithContext(ctx).Where(s.UUID.Eq(uuid)).Take()
}

func (sd *ShareDAO) DeleteByUUID(uuid string) error {
	s := query.Use(db).Share
	_, err := s.WithContext(ctx).Where(s.UUID.Eq(uuid)).Delete()
	return err
}

func (sd *ShareDAO) FindByUserId(userId uint) ([]*model.Share, error) {
	s := query.Use(db).Share
	return s.WithContext(ctx).Preload(s.File).Where(s.UserID.Eq(userId)).Find()
}

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
