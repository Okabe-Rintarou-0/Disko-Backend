package logic

import (
	"context"
	"disko/dao"
	"disko/model"

	"disko/internal/svc"
	"disko/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMySharedFilesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMySharedFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMySharedFilesLogic {
	return &GetMySharedFilesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMySharedFilesLogic) GetMySharedFiles() (resp []*types.ShareDTO, err error) {
	var (
		me     uint
		shares []*model.Share
	)
	me = GetUserId(l.ctx)
	shares, err = l.svcCtx.ShareDAO.FindByUserId(me)
	if err != nil && !dao.IsErrRecordNotFound(err) {
		return nil, err
	}

	var ret []*types.ShareDTO
	for _, s := range shares {
		ret = append(ret, types.FromShare(s, true))
	}

	return ret, nil
}
