package svc

import (
	"cloud_disk/dao"
	"cloud_disk/internal/config"
	"cloud_disk/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config         config.Config
	CheckBlackList rest.Middleware
	UserDAO        dao.IUserDAO
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		CheckBlackList: middleware.NewCheckBlackListMiddleware().Handle,
		UserDAO:        dao.NewUserDAO(),
	}
}
