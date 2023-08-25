package logic

import (
	"context"
	"disko/dao"
	"disko/internal/svc"
	"disko/internal/types"
	"disko/model"
	"fmt"
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
	owner = GetUserId(l.ctx)
	file, err = l.svcCtx.FileDAO.FindById(req.ID)

	if err != nil && !dao.IsErrRecordNotFound(err) {
		return nil, err
	}

	if file == nil {
		return nil, nil
	}

	fmt.Printf("%+v\n", file)

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
