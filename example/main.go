package main

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/vicanso/elton"
	compress "github.com/vicanso/elton-compress"
	"github.com/vicanso/elton/middleware"
)

func main() {
	e := elton.New()
	// 需要注意添加的顺序，选择压缩是按添加的选择顺序选择适合的压缩方式
	// 此处只是示例所有的压缩器，正常使用时，按需使用1，2个压缩方式则可
	config := middleware.NewCompressConfig(
		new(middleware.BrCompressor),
		new(middleware.GzipCompressor),
		new(compress.SnappyCompressor),
		new(compress.ZstdCompressor),
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
		defer resp.Body.Close()
		buf, err := ioutil.ReadAll(resp.Body)
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
