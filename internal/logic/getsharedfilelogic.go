package logic

import (
	"context"
	"disko/dao"
	"disko/model"
	"time"

	"disko/internal/svc"
	"disko/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSharedFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSharedFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSharedFileLogic {
	return &GetSharedFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSharedFileLogic) GetSharedFile(req *types.GetSharedFileRequest) (resp *types.GetSharedFileResponse, err error) {
	var (
		share *model.Share
	)
	share, err = l.svcCtx.ShareDAO.FindByUUIDWithFileAndUser(req.UUID)
	if err != nil && !dao.IsErrRecordNotFound(err) {
		return nil, err
	}

	if share == nil {
		return &types.GetSharedFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "指定的分享文件不存在！",
				Ok:      false,
			},
			Data: nil,
		}, nil
	}

	if share.Password != nil && *share.Password != req.Password {
		return &types.GetSharedFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "提取分享文件的密码错误！",
				Ok:      false,
			},
			Data: nil,
		}, nil
	}

	if share.ExpireAt.Valid && share.ExpireAt.Time.Before(time.Now()) {
		return &types.GetSharedFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "指定的分享文件已过期！",
				Ok:      false,
			},
			Data: nil,
		}, nil
	}

	return &types.GetSharedFileResponse{
		BaseResponse: types.BaseResponse{
			Message: "成功！",
			Ok:      true,
		},
		Data: types.FromShare(share, false),
	}, nil
}
