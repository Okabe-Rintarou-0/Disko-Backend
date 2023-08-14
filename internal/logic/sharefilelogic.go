package logic

import (
	"context"
	"disko/model"
	"errors"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"time"

	"disko/internal/svc"
	"disko/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareFileLogic {
	return &ShareFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareFileLogic) ShareFile(req *types.ShareFileRequest) (resp *types.ShareFileResponse, err error) {
	var (
		file  *model.File
		owner uint
	)

	owner = cast.ToUint(l.ctx.Value("id"))

	file, err = l.svcCtx.FileDAO.FindById(req.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if file == nil {
		return &types.ShareFileResponse{
			Message: "要分享的文件不存在！",
			Ok:      false,
		}, nil
	}

	if file.Owner != owner {
		return &types.ShareFileResponse{
			Message: "无权分享该文件！",
			Ok:      false,
		}, nil
	}

	if file.IsDir {
		return &types.ShareFileResponse{
			Message: "暂不支持分享文件夹！",
			Ok:      false,
		}, nil
	}

	err = l.svcCtx.ShareDAO.Save(&model.Share{
		UUID:     uuid.NewString(),
		Password: req.Password,
		ExpireAt: time.Unix(0, req.ExpireAt*int64(time.Millisecond)),
		FileID:   req.ID,
		UserID:   owner,
	})
	if err != nil {
		return nil, err
	}

	return &types.ShareFileResponse{
		Message: "分享成功！",
		Ok:      false,
	}, nil
}
