package logic

import (
	"context"
	"disko/dao"
	"disko/internal/svc"
	"disko/internal/types"
	"disko/model"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyFilesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyFilesLogic {
	return &GetMyFilesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyFilesLogic) GetMyFiles(req *types.GetMyFileRequest) (resp []*types.FileDTO, err error) {
	var (
		files []*model.File
		dtos  []*types.FileDTO
	)
	ownerId := GetUserId(l.ctx)
	files, err = l.svcCtx.FileDAO.Search(&ownerId, req.Parent, req.Keyword, req.Extensions)

	if err != nil && !dao.IsErrRecordNotFound(err) {
		return nil, err
	}

	dtos = make([]*types.FileDTO, len(files))
	for i, f := range files {
		dtos[i] = types.FromFile(f)
	}

	return dtos, nil
}
