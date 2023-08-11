package logic

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/juju/ratelimit"
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

	// found no file, so just return 404
	if fileMeta == nil {
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
	src := bufio.NewReader(file)

	rate := l.svcCtx.Config.FileStorage.MaxDownloadRate * (1 << 20)
	capacity := cast.ToInt64(rate)
	bucket := ratelimit.NewBucketWithRate(rate, capacity)
	start := time.Now()

	savedName := fileMeta.Name + fileMeta.Ext
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.QueryEscape(savedName)))

	// Copy source to destination, but wrap our reader with rate limited one
	_, err = io.Copy(w, ratelimit.Reader(src, bucket))
	if err != nil {
		return err
	}

	fmt.Printf("Copied in %s\n", time.Since(start))
	return nil
}
