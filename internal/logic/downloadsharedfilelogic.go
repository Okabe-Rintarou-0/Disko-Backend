package logic

import (
	"context"
	"disko/dao"
	"disko/model"
	"net/http"

	"disko/internal/svc"
	"disko/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadSharedFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadSharedFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadSharedFileLogic {
	return &DownloadSharedFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadSharedFileLogic) DownloadSharedFile(req *types.DownloadSharedFileRequest, w http.ResponseWriter) (resp *types.DownloadSharedFileResponse, err error) {
	var (
		share *model.Share
	)

	share, err = l.svcCtx.ShareDAO.FindByUUIDWithFile(req.UUID)
	if err != nil && !dao.IsErrRecordNotFound(err) {
		return nil, err
	}

	if share == nil {
		return &types.DownloadSharedFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "分享不存在！",
				Ok:      false,
			},
		}, nil
	}

	// if password is not matched, return share and a corresponding error
	if share.Password != req.Password {
		return &types.DownloadSharedFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "分享的密码错误！",
				Ok:      false,
			},
		}, nil
	}

	err = downloadFile(w, &share.File, l.svcCtx.Config)
	if err != nil {
		return nil, err
	}

	return &types.DownloadSharedFileResponse{
		BaseResponse: types.BaseResponse{
			Message: "下载成功！",
			Ok:      true,
		},
	}, nil
}
