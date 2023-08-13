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

type GetOneFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOneFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOneFileLogic {
	return &GetOneFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOneFileLogic) GetOneFile(req *types.GetOneFileRequest) (resp *types.FileDTO, err error) {
	var (
		owner uint
		file  *model.File
	)
	owner = cast.ToUint(l.ctx.Value("id"))
	file, err = l.svcCtx.FileDAO.FindById(req.ID)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if file == nil {
		return nil, nil
	}

	// if file is public, or is owned by me
	if !file.Private || file.Owner == owner {
		return &types.FileDTO{
			ID:      file.ID,
			Name:    file.Name,
			Ext:     file.Ext,
			Size:    file.Size,
			UUID:    file.UUID,
			Owner:   file.Owner,
			IsDir:   file.IsDir,
			Private: file.Private,
			Parent:  file.ParentID,
		}, nil
	}

	return nil, nil
}
