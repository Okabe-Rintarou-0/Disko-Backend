package dao

import (
	"disko/model"
	"disko/repository/query"
	"disko/repository/redis"
	"fmt"
)

type IUserDAO interface {
	FindById(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Save(users ...*model.User) error
	UpdateUsage(id uint, delta int64) error
	GetQuotaAndUsage(id uint) (int64, int64, error)
	Logout(token string, expireAt int64) error
}

func NewUserDAO() IUserDAO {
	return &UserDAO{}
}

type UserDAO struct{}

func (d *UserDAO) GetQuotaAndUsage(id uint) (int64, int64, error) {
	u := query.Use(db).User
	user, err := u.WithContext(ctx).Select(u.Quota, u.Usage).Where(u.ID.Eq(id)).Take()
	return user.Quota, user.Usage, err
}

func (d *UserDAO) UpdateUsage(id uint, delta int64) error {
	u := query.Use(db).User
	_, err := u.WithContext(ctx).Where(u.ID.Eq(id)).Update(u.Usage, u.Usage.Add(delta))
	return err
}

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
