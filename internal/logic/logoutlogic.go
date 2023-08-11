package logic

import (
	"context"
	"disko/internal/svc"
	"disko/internal/types"
	"disko/utils"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutRequest) (resp *types.LogoutResponse, err error) {
	var (
		token string
		ok    bool
	)
	token, ok = utils.GetToken(req.Token)
	if !ok {
		return &types.LogoutResponse{
			Message: "登录令牌错误！",
			Ok:      false,
		}, nil
	}

	expireAt := cast.ToInt64(l.ctx.Value("expireAt"))
	if err = l.svcCtx.UserDAO.Logout(token, expireAt); err != nil {
		return nil, err
	}
	return &types.LogoutResponse{
		Message: "登出成功！",
		Ok:      true,
	}, nil
}
