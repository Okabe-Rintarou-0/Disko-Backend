package logic

import (
	"context"
	"disko/constants"
	"disko/dao"
	"disko/internal/svc"
	"disko/internal/types"
	"disko/model"
	"github.com/zeromicro/go-zero/core/logx"
)

type MoveFilesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMoveFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MoveFilesLogic {
	return &MoveFilesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MoveFilesLogic) MoveFile(id, parentID uint) (resp *types.MoveFilesResponse, err error) {
	var (
		owner  uint
		file   *model.File
		parent *model.File
	)
	if id == parentID {
		return &types.MoveFilesResponse{
			BaseResponse: types.BaseResponse{
				Message: "非法操作！请正确指定文件夹！",
				Ok:      false,
			},
		}, nil
	}

	owner = GetUserId(l.ctx)

	file, err = l.svcCtx.FileDAO.FindById(id)
	if err != nil && !dao.IsErrRecordNotFound(err) {
		return nil, err
	}

	if file == nil {
		return &types.MoveFilesResponse{
			BaseResponse: types.BaseResponse{
				Message: "指定的文件不存在！",
				Ok:      false,
			},
		}, nil
	}

	if file.Owner != owner {
		return &types.MoveFilesResponse{
			BaseResponse: types.BaseResponse{
				Message: "非法操作！无权限",
				Ok:      false,
			},
		}, nil
	}

	if parentID != constants.RootDirId {
		parent, err = l.svcCtx.FileDAO.FindById(parentID)
		if err != nil && !dao.IsErrRecordNotFound(err) {
			return nil, err
		}

		if parent == nil {
			return &types.MoveFilesResponse{
				BaseResponse: types.BaseResponse{
					Message: "指定的文件夹不存在！",
					Ok:      false,
				},
			}, nil
		}

		if parent.Owner != owner {
			return &types.MoveFilesResponse{
				BaseResponse: types.BaseResponse{
					Message: "非法操作！无权限",
					Ok:      false,
				},
			}, nil
		}

		if !parent.IsDir || (parent.ParentID != nil && *parent.ParentID == id) {
			return &types.MoveFilesResponse{
				BaseResponse: types.BaseResponse{
					Message: "非法操作！请正确指定文件夹！",
					Ok:      false,
				},
			}, nil
		}

		file.ParentID = &parentID
	} else {
		file.ParentID = nil
	}

	err = l.svcCtx.FileDAO.Save(file)
	if err != nil {
		return nil, err
	}

	return &types.MoveFilesResponse{
		BaseResponse: types.BaseResponse{
			Message: "移动成功！",
			Ok:      true,
		},
	}, nil
}

func (l *MoveFilesLogic) MoveFiles(req *types.MoveFilesRequest) (resp *types.MoveFilesResponse, err error) {
	for _, id := range req.IDs {
		resp, err = l.MoveFile(id, req.Parent)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}
