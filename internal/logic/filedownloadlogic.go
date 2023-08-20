package logic

import (
	"context"
	"disko/dao"
	"disko/internal/svc"
	"disko/internal/types"
	"github.com/spf13/cast"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDownloadLogic {
	return &FileDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileDownloadLogic) FileDownload(req *types.FileDownloadRequest, w http.ResponseWriter) error {
	uuid := req.UUID
	fileMeta, err := l.svcCtx.FileDAO.FindByUUID(uuid)
	if err != nil && !dao.IsErrRecordNotFound(err) {
		return err
	}

	if fileMeta == nil {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	owner := cast.ToUint(l.ctx.Value("id"))
	if fileMeta.Owner != owner && fileMeta.Private {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	// found no file or file is a directory, so just return 404
	if fileMeta == nil || fileMeta.IsDir {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	return downloadFile(w, fileMeta, l.svcCtx.Config)
}
