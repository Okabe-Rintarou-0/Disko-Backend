package logic

import (
	"cloud_disk/dao/model"
	"context"
	"errors"
	"gorm.io/gorm"

	"cloud_disk/internal/svc"
	"cloud_disk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MeLogic {
	return &MeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MeLogic) Me() (resp *types.UserDTO, err error) {
	// todo: add your logic here and delete this line
	var user *model.User
	// todo change id
	user, err = l.svcCtx.UserDAO.FindById(2)
	if user != nil {
		return &types.UserDTO{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return nil, err
}
