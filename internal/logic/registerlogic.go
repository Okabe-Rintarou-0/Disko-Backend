package logic

import (
	"cloud_disk/model"
	"cloud_disk/repository/redis"
	"context"
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"regexp"

	"cloud_disk/internal/svc"
	"cloud_disk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	emailPattern    = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	passwordPattern = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,16}$`
)

var (
	emailRegex    = regexp.MustCompile(emailPattern)
	passwordRegex = regexp2.MustCompile(passwordPattern, 0)
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) verifyEmailFormat(email string) bool {
	return emailRegex.MatchString(email)
}

func (l *RegisterLogic) verifyPasswordFormat(password string) bool {
	succeed, _ := passwordRegex.MatchString(password)
	return succeed
}

func (l *RegisterLogic) encryptPassword(password string) (string, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return cast.ToString(encrypted), nil
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	var (
		user         *model.User
		vcode        string
		encryptedPwd string
	)

	if !l.verifyEmailFormat(req.Email) {
		return &types.RegisterResponse{
			Message: "邮箱格式错误！",
			Ok:      false,
		}, nil
	}

	if !l.verifyPasswordFormat(req.Password) {
		return &types.RegisterResponse{
			Message: "密码格式错误（包含至少一位数字，字母，且长度8-16）！",
			Ok:      false,
		}, nil
	}

	key := fmt.Sprintf("vcode:%s", req.Email)
	vcode, err = redis.Get(key)
	if err != nil {
		return nil, err
	}

	if req.Vcode != vcode {
		return &types.RegisterResponse{
			Message: "验证码错误！",
			Ok:      false,
		}, nil
	}

	user, err = l.svcCtx.UserDAO.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if user != nil {
		return &types.RegisterResponse{
			Message: "注册失败！当前用户已存在！",
			Ok:      false,
		}, nil
	}

	encryptedPwd, err = l.encryptPassword(req.Password)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.UserDAO.Save(&model.User{
		Name:     req.Name,
		Password: encryptedPwd,
		Email:    req.Email,
	})

	if err != nil {
		return nil, err
	}

	// invalidate vcode
	_, err = redis.Del(key)
	if err != nil {
		return nil, err
	}

	return &types.RegisterResponse{
		Message: "注册成功！",
		Ok:      true,
	}, nil
}
