package logic

import (
	"context"
	"disko/model"
	"errors"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"disko/internal/svc"
	"disko/internal/types"

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

func (l *GetMyFilesLogic) GetMyFiles(req *types.GetMyFileRequest) (resp []types.FileDTO, err error) {
	var (
		files []*model.File
		dtos  []types.FileDTO
	)
	ownerId := cast.ToUint(l.ctx.Value("id"))
	files, err = l.svcCtx.FileDAO.Search(&ownerId, req.Parent, req.Keyword, req.Extensions)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	dtos = make([]types.FileDTO, len(files))
	for i, f := range files {
		dtos[i] = types.FileDTO{
			ID:      f.ID,
			Name:    f.Name,
			Ext:     f.Ext,
			Size:    f.Size,
			UUID:    f.UUID,
			Owner:   f.Owner,
			IsDir:   f.IsDir,
			Private: f.Private,
			Parent:  f.ParentID,
		}
	}

	return dtos, nil
}
