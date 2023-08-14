package logic

import (
	"context"
	"disko/constants"
	"disko/model"
	"errors"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"path"

	"disko/internal/svc"
	"disko/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFileLogic {
	return &UpdateFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFileLogic) UpdateFile(req *types.UpdateFileRequest) (resp *types.UpdateFileResponse, err error) {
	var (
		file   *model.File
		parent *model.File
	)

	file, err = l.svcCtx.FileDAO.FindById(req.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if file == nil {
		return &types.UpdateFileResponse{
			Message: "指定的文件不存在！",
			Ok:      false,
		}, nil
	}

	owner := cast.ToUint(l.ctx.Value("id"))
	if file.Owner != owner {
		return &types.UpdateFileResponse{
			Message: "非法操作！无权限！",
			Ok:      false,
		}, nil
	}

	if req.Parent != nil && *req.Parent > 0 {
		// cannot move under itself
		if req.ID == *req.Parent {
			return &types.UpdateFileResponse{
				Message: "非法操作！请正确指定文件夹！",
				Ok:      false,
			}, nil
		}

		parent, err = l.svcCtx.FileDAO.FindById(*req.Parent)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		if parent == nil {
			return &types.UpdateFileResponse{
				Message: "指定的文件夹不存在！",
				Ok:      false,
			}, nil
		}

		// avoid circle
		if !parent.IsDir || (parent.ParentID != nil && *parent.ParentID == req.ID) {
			return &types.UpdateFileResponse{
				Message: "非法操作！请正确指定文件夹！",
				Ok:      false,
			}, nil
		}

		if parent.Owner != owner {
			return &types.UpdateFileResponse{
				Message: "非法操作！无权限！",
				Ok:      false,
			}, nil
		}

		// use zero to refer null
		file.ParentID = req.Parent
	}

	if req.Parent != nil && *req.Parent == constants.RootDirId {
		file.ParentID = nil
	}

	if len(req.Name) > 0 {
		// name with ext
		ext := path.Ext(req.Name)
		file.Ext = ext
		file.Name = req.Name[:len(req.Name)-len(ext)]
	}

	if req.Private != nil {
		file.Private = *req.Private
	}

	err = l.svcCtx.FileDAO.Save(file)
	if err != nil {
		return nil, err
	}

	return &types.UpdateFileResponse{
		Message: "更新成功！",
		Ok:      true,
	}, nil
}
