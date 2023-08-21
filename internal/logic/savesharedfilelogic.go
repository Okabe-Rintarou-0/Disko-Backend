package logic

import (
	"context"
	"disko/dao"
	"disko/model"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"path"
	"time"

	"disko/internal/svc"
	"disko/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveSharedFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveSharedFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveSharedFileLogic {
	return &SaveSharedFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveSharedFileLogic) SaveSharedFile(req *types.SaveSharedFileRequest) (resp *types.SaveSharedFileResponse, err error) {
	var (
		share    *model.Share
		owner    uint
		filename string
		ext      string
		existed  *model.File
		quota    int64
		usage    int64
	)

	share, err = l.svcCtx.ShareDAO.FindByUUIDWithFile(req.UUID)
	if err != nil && !dao.IsErrRecordNotFound(err) {
		return nil, err
	}

	if share == nil {
		return &types.SaveSharedFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "分享不存在！",
				Ok:      false,
			},
		}, nil
	}

	if share.ExpireAt.Valid && share.ExpireAt.Time.Before(time.Now()) {
		return &types.SaveSharedFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "分享已过期！",
				Ok:      false,
			},
		}, nil
	}

	// if password is not matched, return share and a corresponding error
	if req.Password != nil && share.Password != nil && *share.Password != *req.Password {
		return &types.SaveSharedFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "分享的密码错误！",
				Ok:      false,
			},
		}, nil
	}

	owner = cast.ToUint(l.ctx.Value("id"))

	if len(req.Name) > 0 {
		ext = path.Ext(req.Name)
		filename = req.Name[:len(req.Name)-len(ext)]
	} else {
		ext = share.File.Ext
		filename = share.File.Name
	}

	existed, err = l.svcCtx.FileDAO.FindByOwnerAndParentAndName(owner, nil, filename)
	if err != nil && !dao.IsErrRecordNotFound(err) {
		return nil, err
	}

	if existed != nil && !existed.IsDir {
		return &types.SaveSharedFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "已存在同名文件！",
				Ok:      false,
			},
		}, nil
	}

	uid := uuid.NewString()
	quota, usage, err = l.svcCtx.UserDAO.GetQuotaAndUsage(owner)
	if err != nil {
		return nil, err
	}

	if quota < usage+share.File.Size {
		return &types.SaveSharedFileResponse{
			BaseResponse: types.BaseResponse{
				Message: "空间不足！",
				Ok:      false,
			},
		}, nil
	}

	// it must be a file
	// since directory can only be created instead of uploaded
	err = l.svcCtx.FileDAO.Save(&model.File{
		Name: filename,
		Ext:  ext,
		Size: share.File.Size,
		UUID: uid,
		// multiple references
		Path:  share.File.Path,
		Owner: owner,
		IsDir: false,
		// private file by default
		Private:  true,
		ParentID: nil,
		Parent:   nil,
	})
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.UserDAO.UpdateUsage(owner, share.File.Size)
	if err != nil {
		return nil, err
	}

	return &types.SaveSharedFileResponse{
		BaseResponse: types.BaseResponse{
			Message: "保存成功！",
			Ok:      true,
		},
	}, nil
}
