package logic

import (
	"context"
	"database/sql"
	"disko/dao"
	"disko/model"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"regexp"
	"time"

	"disko/internal/svc"
	"disko/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var sharePasswordRegex = regexp.MustCompile("[0-9a-zA-Z]{4}")

type ShareFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareFileLogic {
	return &ShareFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareFileLogic) ShareFile(req *types.ShareFileRequest) (resp *types.ShareFileResponse, err error) {
	var (
		file  *model.File
		owner uint
	)

	owner = cast.ToUint(l.ctx.Value("id"))

	file, err = l.svcCtx.FileDAO.FindById(req.ID)
	if err != nil && !dao.IsErrRecordNotFound(err) {
		return nil, err
	}

	if file == nil {
		return &types.ShareFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "要分享的文件不存在！",
				Ok:      false,
			},
		}, nil
	}

	if file.Owner != owner {
		return &types.ShareFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "无权分享该文件！",
				Ok:      false,
			},
		}, nil
	}

	if file.IsDir {
		return &types.ShareFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "暂不支持分享文件夹！",
				Ok:      false,
			},
		}, nil
	}

	if req.Password != nil && !sharePasswordRegex.MatchString(*req.Password) {
		return &types.ShareFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "密码格式错误！",
				Ok:      false,
			},
		}, nil
	}

	var expireAt sql.NullTime
	if req.ExpireAt != nil {
		expireAt = sql.NullTime{
			Time:  time.Unix(0, *req.ExpireAt*int64(time.Millisecond)),
			Valid: true,
		}
	}

	err = l.svcCtx.ShareDAO.Save(&model.Share{
		UUID:     uuid.NewString(),
		Password: req.Password,
		ExpireAt: expireAt,
		FileID:   req.ID,
		UserID:   owner,
	})
	if err != nil {
		return nil, err
	}

	return &types.ShareFileResponse{
		BaseResponse: types.BaseResponse{
			Message: "分享成功！",
			Ok:      true,
		},
	}, nil
}
