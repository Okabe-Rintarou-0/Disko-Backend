package logic

import (
	"context"
	"disko/model"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"disko/internal/svc"
	"disko/internal/types"

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
	id := cast.ToUint(l.ctx.Value("id"))
	user, err = l.svcCtx.UserDAO.FindById(id)
	fmt.Printf("id = %d, me: %+v", id, user)
	if user != nil {
		return &types.UserDTO{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Quota: user.Quota,
			Usage: user.Usage,
		}, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return nil, err
}
