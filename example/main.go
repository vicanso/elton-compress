package main

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/vicanso/elton"
	compress "github.com/vicanso/elton-compress"
)

func main() {
	e := elton.New()
	compressConfig := compress.Config{
		// 大于1KB的数据才做压缩
		MinLength: 1024,
		Levels: map[string]int{
			compress.GzipEncoding: 6,
			compress.BrEncoding:   6,
			compress.Lz4Encoding:  6,
			compress.S2Encoding:   2,
			compress.ZstdEncoding: 2,
		},
	}
	// 需要注意添加的顺序，选择压缩是按添加的选择顺序选择适合的压缩方式
	compressConfig.AddCompressor(new(compress.BrCompressor))
	compressConfig.AddCompressor(new(compress.GzipCompressor))
	compressConfig.AddCompressor(new(compress.SnappyCompressor))
	compressConfig.AddCompressor(new(compress.ZstdCompressor))
	compressConfig.AddCompressor(new(compress.S2Compressor))
	compressConfig.AddCompressor(new(compress.Lz4Compressor))

	e.Use(compress.New(compressConfig))

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
	e.ListenAndServe(":3000")
}
