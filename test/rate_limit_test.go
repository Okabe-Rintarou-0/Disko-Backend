package test

import (
	"fmt"
	"github.com/juju/ratelimit"
	"github.com/spf13/cast"
	"io"
	"os"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	file, err := os.Open("./test.mp4")
	if err != nil {
		panic(err)
	}
	rate := cast.ToFloat64(10 * (1 << 20))
	capacity := cast.ToInt64(rate)
	bucket := ratelimit.NewBucketWithRate(rate, capacity)
	start := time.Now()

	tempFile, err := os.Create("tmp.mp4")
	if err != nil {
		panic(err)
	}

	go func() {
		for range time.Tick(time.Millisecond * 100) {
			file, err := os.Open("./tmp.mp4")
			if err != nil {
				fmt.Println("无法打开文件:", err)
				return
			}
			defer file.Close()

			// 获取文件状态信息
			fileInfo, err := file.Stat()
			if err != nil {
				fmt.Println("无法获取文件信息:", err)
				return
			}
			fileSize := fileInfo.Size()
			fmt.Printf("文件大小: %d 字节\n", fileSize)
		}
	}()

	// Copy source to destination, but wrap our reader with rate limited one
	_, err = io.Copy(tempFile, ratelimit.Reader(file, bucket))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Copied in %s\n", time.Since(start))
}
