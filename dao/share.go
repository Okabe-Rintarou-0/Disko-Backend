package dao

import (
	"disko/model"
	"disko/repository/query"
)

type IShareDAO interface {
	FindByID(id uint) (*model.Share, error)
	Save(shares ...*model.Share) error
}

func NewShareDAO() IShareDAO {
	return &ShareDAO{}
}

type ShareDAO struct{}

func (sd *ShareDAO) FindByID(id uint) (*model.Share, error) {
	s := query.Use(db).Share
	return s.WithContext(ctx).Where(s.ID.Eq(id)).Take()
}

func (sd *ShareDAO) Save(shares ...*model.Share) error {
	s := query.Use(db).Share
	return s.WithContext(ctx).Save(shares...)
}
