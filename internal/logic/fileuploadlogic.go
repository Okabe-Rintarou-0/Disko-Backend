package logic

import (
	"context"
	"disko/dao"
	"disko/internal/svc"
	"disko/internal/types"
	"disko/model"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/zeromicro/go-zero/core/logx"
)

const bytesPerGB = 1 << 30

type FileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadLogic) FileUpload(r *http.Request, req *types.FileUploadRequest) (resp *types.FileUploadResponse, err error) {
	var (
		file     multipart.File
		header   *multipart.FileHeader
		tempFile *os.File
		uid      string
		filePath string
		existed  *model.File
		parent   *model.File
		quota    int64
		usage    int64
	)
	owner := GetUserId(l.ctx)
	// convert to bytes
	// 1 GB = (1 << 30) B
	maxFileSize := cast.ToInt64(l.svcCtx.Config.FileStorage.MaxFileSize * bytesPerGB)
	err = r.ParseMultipartForm(maxFileSize)
	if err != nil {
		return nil, err
	}
	file, header, err = r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if header.Size > maxFileSize {
		return &types.FileUploadResponse{
			BaseResponse: types.BaseResponse{
				Message: fmt.Sprintf("超过文件大小上限：%fGB", l.svcCtx.Config.FileStorage.MaxFileSize),
				Ok:      false,
			},
		}, nil
	}

	if req.Parent != nil {
		parent, err = l.svcCtx.FileDAO.FindById(*req.Parent)
		if err != nil && !dao.IsErrRecordNotFound(err) {
			return nil, err
		}

		// specified parent does not exist
		if parent == nil {
			return &types.FileUploadResponse{
				BaseResponse: types.BaseResponse{
					Message: "指定的文件夹不存在！",
					Ok:      false,
				},
			}, nil
		}

		// if parent does not belong to me, then I have no authority to create a file under it
		if parent.Owner != owner {
			return &types.FileUploadResponse{
				BaseResponse: types.BaseResponse{
					Message: "非法操作！无权限！",
					Ok:      false,
				},
			}, nil
		}
	}

	ext := path.Ext(header.Filename)
	filename := header.Filename[:len(header.Filename)-len(ext)]
	fmt.Printf("filename = %s, ext = %s\n", filename, ext)

	// only one case is invalid: same name under same parent directory
	existed, err = l.svcCtx.FileDAO.FindByOwnerAndParentAndName(owner, req.Parent, filename)
	if err != nil && !dao.IsErrRecordNotFound(err) {
		return nil, err
	}

	// if ABC is a file(wo extension), while another ABC is a directory, it is ok
	// this statement means: 'if exists file with same name'
	if existed != nil && !existed.IsDir {
		return &types.FileUploadResponse{
			BaseResponse: types.BaseResponse{
				Message: "已存在同名文件！",
				Ok:      false,
			},
		}, nil
	}

	dir := l.svcCtx.Config.FileStorage.LocalStorage.Dir
	// save as <uuid><ext>
	uid = uuid.NewString()
	filePath = path.Join(dir, uid+ext)

	tempFile, err = os.Create(filePath)
	if err != nil {
		return nil, err
	}

	defer tempFile.Close()
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return nil, err
	}

	quota, usage, err = l.svcCtx.UserDAO.GetQuotaAndUsage(owner)
	if err != nil {
		return nil, err
	}

	if quota < usage+header.Size {
		return &types.FileUploadResponse{
			BaseResponse: types.BaseResponse{
				Message: "空间不足！",
				Ok:      false,
			},
		}, nil
	}

	// it must be a file
	// since directory can only be created instead of uploaded
	err = l.svcCtx.FileDAO.Save(&model.File{
		Name:  filename,
		Ext:   ext,
		Size:  header.Size,
		UUID:  uid,
		Path:  uid + ext,
		Owner: owner,
		IsDir: false,
		// private file by default
		Private:  true,
		ParentID: req.Parent,
		Parent:   parent,
	})
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.UserDAO.UpdateUsage(owner, header.Size)
	if err != nil {
		return nil, err
	}

	return &types.FileUploadResponse{
		BaseResponse: types.BaseResponse{
			Message: "上传成功！",
			Ok:      true,
		},
	}, nil
}
