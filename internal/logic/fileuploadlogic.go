package logic

import (
	"context"
	"disko/internal/svc"
	"disko/internal/types"
	"disko/model"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"gorm.io/gorm"
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
		parent   *model.File
		parentID *uint
	)
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

	fmt.Printf("Uploaded File: %+v\n", header.Filename)
	fmt.Printf("File Size: %+v\n", header.Size)
	fmt.Printf("MIME Header: %+v\n", header.Header)

	if header.Size > maxFileSize {
		return &types.FileUploadResponse{
			Message: fmt.Sprintf("超过文件大小上限：%fGB", l.svcCtx.Config.FileStorage.MaxFileSize),
			Ok:      false,
		}, nil
	}

	if len(req.Parent) > 0 {
		parent, err = l.svcCtx.FileDAO.FindByUUID(req.Parent)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return &types.FileUploadResponse{
			Message: "指定的文件目录不存在！",
			Ok:      false,
		}, nil
	}

	if parent != nil {
		parentID = &parent.ID
	}

	ext := path.Ext(header.Filename)
	filename := header.Filename[:len(header.Filename)-len(ext)]
	fmt.Printf("filename = %s, ext = %s\n", filename, ext)

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

	// it must be a file
	// since directory can only be created instead of uploaded
	err = l.svcCtx.FileDAO.Save(&model.File{
		Name:     filename,
		Ext:      ext,
		Size:     header.Size,
		UUID:     uid,
		Path:     uid + ext,
		Owner:    cast.ToUint(l.ctx.Value("id")),
		IsDir:    false,
		ParentID: parentID,
		Parent:   parent,
	})
	if err != nil {
		return nil, err
	}

	return &types.FileUploadResponse{
		Message: "上传成功！",
		Ok:      true,
	}, nil
}
