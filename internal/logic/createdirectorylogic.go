package logic

import (
	"context"
	"disko/dao"
	"disko/model"
	"github.com/google/uuid"
	"path"

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

	if len(req.Name) == 0 {
		return &types.CreateDirectoryResponse{
			BaseResponse: types.BaseResponse{
				Message: "文件夹名字不得为空！",
				Ok:      false,
			},
		}, nil
	}

	owner := GetUserId(l.ctx)

	// step 1 check whether parent exists
	if req.Parent != nil {
		parent, err = l.svcCtx.FileDAO.FindById(*req.Parent)
		if err != nil && !dao.IsErrRecordNotFound(err) {
			return nil, err
		}

		if parent == nil {
			return &types.CreateDirectoryResponse{
				BaseResponse: types.BaseResponse{
					Message: "指定的文件夹不存在！",
					Ok:      false,
				},
			}, nil
		}

		// if parent does not belong to me, then I have no authority to create a file under it
		if parent.Owner != owner {
			return &types.CreateDirectoryResponse{
				BaseResponse: types.BaseResponse{
					Message: "非法操作！无权限！",
					Ok:      false,
				},
			}, nil
		}
	}

	existed, err = l.svcCtx.FileDAO.FindByOwnerAndParentAndName(owner, req.Parent, req.Name)
	if err != nil && !dao.IsErrRecordNotFound(err) {
		return nil, err
	}

	if existed != nil && existed.IsDir {
		return &types.CreateDirectoryResponse{
			BaseResponse: types.BaseResponse{
				Message: "已存在同名文件夹！",
				Ok:      false,
			},
		}, nil
	}
	parentPath := ""
	if parent != nil {
		parentPath = parent.Name
	}

	// dir is logical, we don't need to create a real one
	err = l.svcCtx.FileDAO.Save(&model.File{
		Name:     path.Join(parentPath, req.Name),
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
		BaseResponse: types.BaseResponse{
			Message: "创建成功！",
			Ok:      true,
		},
	}, nil
}
