package logic

import (
	"cloud_disk/dao/model"
	"context"
	"errors"
	"github.com/spf13/cast"
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
	var user *model.User
	var id = cast.ToUint(l.ctx.Value("id"))
	user, err = l.svcCtx.UserDAO.FindById(id)
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
