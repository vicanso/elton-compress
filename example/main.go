package main

import (
	"bytes"
	"io"
	"net/http"

	compress "github.com/vicanso/elton-compress/v2"
	"github.com/vicanso/elton/v2"
	"github.com/vicanso/elton/v2/middleware"
)

func main() {
	e := elton.New()
	// 需要注意添加的顺序，选择压缩是按添加的选择顺序选择适合的压缩方式
	// 此处只是示例所有的压缩器，正常使用时，按需使用1，2个压缩方式则可
	config := middleware.NewCompressConfig(
		// gzip、br、zstd 由 elton 内置
		middleware.NewBrCompressor(),
		middleware.NewGzipCompressor(),
		middleware.NewZstdCompressor(),
		// snappy、s2、lz4 由本模块提供
		compress.NewSnappyCompressor(),
		compress.NewS2Compressor(),
		&compress.Lz4Compressor{
			MinLength: 10 * 1024,
		},
	)
	e.Use(middleware.NewCompress(config))

	e.GET("/", func(c *elton.Context) (err error) {
		resp, err := http.Get("https://code.jquery.com/jquery-3.4.1.min.js")
		if err != nil {
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		c.SetContentTypeByExt(".js")
		c.BodyBuffer = bytes.NewBuffer(buf)
		return
	})
	err := e.ListenAndServe(":3000")
	if err != nil {
		panic(err)
	}
}
