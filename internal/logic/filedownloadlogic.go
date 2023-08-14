package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"disko/internal/svc"
	"disko/internal/types"

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
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
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

	var filePath string
	// remote url?
	if strings.HasPrefix(fileMeta.Path, "http") {
		filePath = fileMeta.Path
	} else {
		// local url
		filePath = path.Join(l.svcCtx.Config.FileStorage.LocalStorage.Dir, fileMeta.Path)
	}

	var file *os.File
	file, err = os.Open(filePath)
	//src := bufio.NewReader(file)

	rate := l.svcCtx.Config.FileStorage.MaxDownloadRate * (1 << 20)
	//capacity := cast.ToInt64(rate)
	//bucket := ratelimit.NewBucketWithRate(rate, capacity)

	savedName := fileMeta.Name + fileMeta.Ext
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.QueryEscape(savedName)))
	w.Header().Set("Content-Length", cast.ToString(fileMeta.Size))

	chunkNum := cast.ToInt(cast.ToFloat64(fileMeta.Size) / rate)
	chunkSize := cast.ToInt64(rate)
	for i := 0; i < chunkNum; i++ {
		_, err = io.CopyN(w, file, chunkSize)
		if err != nil {
			return err
		}
		w.(http.Flusher).Flush()
		time.Sleep(time.Second)
	}

	chunkSize = fileMeta.Size % chunkSize
	if chunkSize > 0 {
		_, err = io.CopyN(w, file, chunkSize)
		if err != nil {
			return err
		}
		w.(http.Flusher).Flush()
		time.Sleep(time.Second)
	}

	// Copy source to destination, but wrap our reader with rate limited one
	//_, err = io.Copy(w, ratelimit.Reader(src, bucket))
	//if err != nil {
	//	return err
	//}
	return nil
}
