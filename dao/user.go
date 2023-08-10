package dao

import (
	"cloud_disk/model"
	"cloud_disk/repository/query"
	"cloud_disk/repository/redis"
	"fmt"
)

type IUserDAO interface {
	FindById(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Save(users ...*model.User) error
	Logout(token string, expireAt int64) error
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

func (d *UserDAO) Logout(token string, expireAt int64) error {
	var (
		key string
		err error
	)

	key = fmt.Sprintf("blacklist:%s", token)

	if err = redis.Set(key, ""); err != nil {
		return err
	}

	return redis.ExpireAt(key, expireAt)
}
