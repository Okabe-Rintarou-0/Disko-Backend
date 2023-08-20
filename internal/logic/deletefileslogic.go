package logic

import (
	"context"

	"disko/internal/svc"
	"disko/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFilesLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	deleteFileLogic *DeleteFileLogic
}

func NewDeleteFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFilesLogic {
	return &DeleteFilesLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		deleteFileLogic: NewDeleteFileLogic(ctx, svcCtx),
	}
}

func (l *DeleteFilesLogic) DeleteFiles(req *types.DeleteFilesRequest) (resp *types.DeleteFilesResponse, err error) {
	var (
		res *types.DeleteFileResponse
	)
	for _, uuid := range req.UUIDs {
		res, err = l.deleteFileLogic.DeleteFile(&types.DeleteFileRequest{UUID: uuid})
		if err != nil {
			return nil, err
		}
		if !res.Ok {
			break
		}
	}

	return &types.DeleteFilesResponse{
		BaseResponse: types.BaseResponse{
			Message: res.Message,
			Ok:      res.Ok,
		},
	}, nil
}
