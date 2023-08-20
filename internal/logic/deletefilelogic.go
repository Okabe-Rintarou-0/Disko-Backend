package logic

import (
	"context"
	"disko/dao"
	"disko/internal/svc"
	"disko/internal/types"
	"disko/model"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFileLogic) DeleteFile(req *types.DeleteFileRequest) (resp *types.DeleteFileResponse, err error) {
	var (
		file  *model.File
		owner uint
	)

	owner = cast.ToUint(l.ctx.Value("id"))
	file, err = l.svcCtx.FileDAO.FindByUUID(req.UUID)
	if err != nil && !dao.IsErrRecordNotFound(err) {
		return nil, err
	}

	if file == nil {
		return &types.DeleteFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "指定的文件或文件夹不存在！",
				Ok:      false,
			},
		}, nil
	}

	if file.Owner != owner {
		return &types.DeleteFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "非法操作！无权限！",
				Ok:      false,
			},
		}, nil
	}

	err = l.svcCtx.FileDAO.DeleteByUUID(req.UUID)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.UserDAO.UpdateUsage(owner, -file.Size)
	if err != nil {
		return nil, err
	}

	// check shared files

	return &types.DeleteFileResponse{
		BaseResponse: types.BaseResponse{
			Message: "删除成功！",
			Ok:      true,
		},
	}, nil
}
