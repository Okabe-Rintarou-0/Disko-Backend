package logic

import (
	"context"
	"disko/email"
	"disko/repository/redis"
	"fmt"
	"github.com/adjust/rmq/v5"
	"github.com/spf13/cast"
	"math/rand"
	"time"

	"disko/internal/svc"
	"disko/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const template = "您好！<br>请使用以下验证码完成登录，10分钟内有效：<br>%d<br>如非本人操作请忽略该邮件。<br>Disko"

var (
	mq     rmq.Queue
	ctx    = context.Background()
	logger = logx.WithContext(ctx)
)

func init() {
	rand.Seed(time.Now().UnixNano())
	errChan := make(chan error)

	connection, err := rmq.OpenConnection("register", "tcp", redis.Host(), 1, errChan)
	if err != nil {
		panic(err)
	}
	go func() {
		for err = range errChan {
			logger.Error(err)
		}
	}()
	mq, err = connection.OpenQueue("vcode")
	if err != nil {
		panic(err)
	}

	err = mq.StartConsuming(10, time.Second)
	if err != nil {
		panic(err)
	}

	c := consumer{}
	_, err = mq.AddConsumer("task-consumer", &c)
	if err != nil {
		panic(err)
	}
}

type consumer struct{}

func (c *consumer) Consume(delivery rmq.Delivery) {
	var (
		err error
		key string
	)
	to := delivery.Payload()
	vcode := c.genRandomVcode()
	logger.Infof("准备发送验证码 %d 到 %s", vcode, to)

	if err = c.sendVcode(to, vcode); err != nil {
		logger.Error(err)
		goto ack
	}
	logger.Infof("发送验证码成功！")

	// record vcode
	key = fmt.Sprintf("vcode:%s", to)
	if err = redis.Set(key, cast.ToString(vcode)); err != nil {
		logger.Error(err)
	}
	// valid for 10 minutes
	if err = redis.Expire(key, 600); err != nil {
		logger.Error(err)
	}
ack:
	if err = delivery.Ack(); err != nil {
		// handle ack error
		logger.Error(err)
	}
}

func (c *consumer) genRandomVcode() int {
	return rand.Intn(9000) + 1000
}

func (c *consumer) sendVcode(to string, vcode int) error {
	content := fmt.Sprintf(template, vcode)
	return email.SendEmail(to, "Disko 注册验证码", content)
}

type RegisterVcodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterVcodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterVcodeLogic {
	return &RegisterVcodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterVcodeLogic) RegisterVcode(req *types.RegisterVcodeRequest) (resp *types.RegisterVcodeResponse, err error) {
	// send email vcode task to mq
	// sending email takes some time, we don't need to hang the http request
	// async is a good way
	err = mq.Publish(req.Email)
	if err == nil {
		return &types.RegisterVcodeResponse{
			BaseResponse: types.BaseResponse{
				Message: "发送验证码成功！",
				Ok:      true,
			},
		}, nil
	}
	return nil, err
}
