package svc

import (
	"disko/dao"
	"disko/internal/config"
	"disko/internal/middleware"
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
