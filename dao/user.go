package dao

import (
	"cloud_disk/dao/model"
	"cloud_disk/dao/query"
)

type IUserDAO interface {
	FindById(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Save(users ...*model.User) error
}

func NewUserDAO() IUserDAO {
	return &UserDAO{}
}

type UserDAO struct{}

func (d *UserDAO) FindById(id uint) (*model.User, error) {
	u := query.Use(db).User
	return u.WithContext(ctx).Where(u.ID.Eq(id)).Take()
}

func (d *UserDAO) FindByEmail(email string) (*model.User, error) {
	u := query.Use(db).User
	return u.WithContext(ctx).Where(u.Email.Eq(email)).Take()
}

func (d *UserDAO) Save(users ...*model.User) error {
	u := query.Use(db).User
	return u.WithContext(ctx).Save(users...)
}
