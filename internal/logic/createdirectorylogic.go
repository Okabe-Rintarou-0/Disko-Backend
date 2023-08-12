package logic

import (
	"context"
	"disko/model"
	"errors"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"disko/internal/svc"
	"disko/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDirectoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDirectoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDirectoryLogic {
	return &CreateDirectoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDirectoryLogic) CreateDirectory(req *types.CreateDirectoryRequest) (resp *types.CreateDirectoryResponse, err error) {
	var (
		parent  *model.File
		existed *model.File
	)
	owner := cast.ToUint(l.ctx.Value("id"))

	// step 1 check whether parent exists
	if req.Parent != nil {
		parent, err = l.svcCtx.FileDAO.FindById(*req.Parent)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		if parent == nil {
			return &types.CreateDirectoryResponse{
				Message: "指定的文件夹不存在！",
				Ok:      false,
			}, nil
		}

		// if parent does not belong to me, then I have no authority to create a file under it
		if parent.Owner != owner {
			return &types.CreateDirectoryResponse{
				Message: "非法操作！无权限！",
				Ok:      false,
			}, nil
		}
	}

	existed, err = l.svcCtx.FileDAO.FindByOwnerAndParentAndName(owner, req.Parent, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existed != nil && existed.IsDir {
		return &types.CreateDirectoryResponse{
			Message: "已存在同名文件夹！",
			Ok:      false,
		}, nil
	}

	// dir is logical, we don't need to create a real one
	err = l.svcCtx.FileDAO.Save(&model.File{
		Name:     req.Name,
		Ext:      "",
		Size:     0,
		UUID:     uuid.NewString(),
		Path:     "",
		Owner:    owner,
		IsDir:    true,
		Private:  true,
		ParentID: req.Parent,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateDirectoryResponse{
		Message: "创建成功！",
		Ok:      true,
	}, nil
}
