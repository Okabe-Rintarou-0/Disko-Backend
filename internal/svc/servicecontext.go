package svc

import (
	"cloud_disk/dao"
	"cloud_disk/internal/config"
)

type ServiceContext struct {
	Config  config.Config
	UserDAO dao.IUserDAO
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserDAO: dao.NewUserDAO(),
	}
}
