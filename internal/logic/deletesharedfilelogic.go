package logic

import (
	"context"
	"disko/dao"
	"disko/model"

	"disko/internal/svc"
	"disko/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSharedFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSharedFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSharedFileLogic {
	return &DeleteSharedFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSharedFileLogic) DeleteSharedFile(req *types.DeleteSharedFileRequest) (resp *types.DeleteSharedFileResponse, err error) {
	var (
		share *model.Share
		owner uint
	)
	share, err = l.svcCtx.ShareDAO.FindByUUID(req.UUID)
	if err != nil && !dao.IsErrRecordNotFound(err) {
		return nil, err
	}

	if share == nil {
		return &types.DeleteSharedFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "删除失败！指定的分享不存在！",
				Ok:      false,
			},
		}, nil
	}

	owner = GetUserId(l.ctx)
	if owner != share.UserID {
		return &types.DeleteSharedFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "删除失败！无权限！",
				Ok:      false,
			},
		}, nil
	}

	err = l.svcCtx.ShareDAO.DeleteByUUID(req.UUID)
	if err != nil {
		return nil, err
	}

	return &types.DeleteSharedFileResponse{
		BaseResponse: types.BaseResponse{
			Message: "删除成功！",
			Ok:      true,
		},
	}, nil
}
