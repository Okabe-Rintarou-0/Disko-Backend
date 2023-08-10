package logic

import (
	"cloud_disk/dao/model"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"time"

	"cloud_disk/internal/svc"
	"cloud_disk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	var (
		user     *model.User
		jwtToken string
	)
	l.Logger.Infof("Receive email: %s and password: %s", req.Email, req.Password)

	user, err = l.svcCtx.UserDAO.FindByEmail(req.Email)
	if user == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Yes")
			err = nil
		}
		return &types.LoginResponse{
			Message: "邮箱或密码错误！",
			Ok:      false,
		}, err
	}

	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessSecret := l.svcCtx.Config.Auth.AccessSecret
	jwtToken, err = l.getJwtToken(accessSecret, now, accessExpire, user.ID)

	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		Message:  "登录成功！",
		Ok:       true,
		Token:    jwtToken,
		ExpireAt: now + accessExpire,
	}, err
}

func (l *LoginLogic) getJwtToken(secret string, nowDate, accessExpire int64, id uint) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = nowDate + accessExpire
	claims["iat"] = nowDate
	claims["id"] = id
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secret))
}
