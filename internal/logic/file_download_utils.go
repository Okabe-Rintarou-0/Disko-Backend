package logic

import (
	"disko/internal/config"
	"disko/model"
	"fmt"
	"github.com/spf13/cast"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

func downloadFile(w http.ResponseWriter, fileMeta *model.File, config config.Config) error {
	var (
		filePath string
		err      error
		file     *os.File
	)
	// remote url?
	if strings.HasPrefix(fileMeta.Path, "http") {
		filePath = fileMeta.Path
	} else {
		// local url
		filePath = path.Join(config.FileStorage.LocalStorage.Dir, fileMeta.Path)
	}

	file, err = os.Open(filePath)
	//src := bufio.NewReader(file)

	rate := config.FileStorage.MaxDownloadRate * (1 << 20)
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
