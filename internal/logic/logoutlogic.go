package logic

import (
	"cloud_disk/internal/svc"
	"cloud_disk/internal/types"
	"cloud_disk/utils"
	"github.com/spf13/cast"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	req    *http.Request
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(req *http.Request, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(req.Context()),
		req:    req,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout() (resp *types.LogoutResponse, err error) {
	var (
		token string
		ok    bool
	)
	token, ok = utils.GetToken(l.req)
	if !ok {
		return &types.LogoutResponse{
			Message: "登录令牌错误！",
			Ok:      false,
		}, nil
	}

	expireAt := cast.ToInt64(l.req.Context().Value("expireAt"))
	if err = l.svcCtx.UserDAO.Logout(token, expireAt); err != nil {
		return nil, err
	}
	return &types.LogoutResponse{
		Message: "登出成功！",
		Ok:      false,
	}, nil
}
