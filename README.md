# elton-compress

[![license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/vicanso/elton-compress/blob/master/LICENSE)
[![Build Status](https://github.com/vicanso/elton-compress/workflows/Test/badge.svg)](https://github.com/vicanso/elton-compress/actions)

为 [elton](https://github.com/vicanso/elton) 的 compress 中间件补充 elton 未内置的压缩算法。适用于对 CPU 敏感、压缩率要求不高的内网场景。

```bash
go get github.com/vicanso/elton-compress/v2
```

从 1.x 升级请参考 [2.0 迁移指南](./docs/migration-v2.md)。

## 压缩器

| 构造函数 | 编码 | 级别 | 说明 |
|---|---|---|---|
| `NewSnappyCompressor()` | `snz` | 1-3，默认 2 | snappy block 格式，自带原始长度前缀，标准 snappy 解码器可直接解压 |
| `NewS2Compressor()` | `s2` | 1-3，默认 1 | s2 block 格式，是 snappy 的超集，速度与压缩率均优于 snappy，但只能由 s2 解码器解压 |
| `NewLz4Compressor()` | `lz4` | 1-10，默认 1 | 标准 LZ4 frame 格式，任意 lz4 客户端可直接解压 |

gzip、br、zstd 已由 `github.com/vicanso/elton/v2/middleware` 内置，请直接使用 `NewGzipCompressor()`、`NewBrCompressor()`、`NewZstdCompressor()`。

每个压缩器都提供对应的包级函数，可脱离中间件单独使用：

```go
buf, err := compress.Lz4Compress(data, compress.Lz4LevelBest)
raw, err := compress.Lz4Decompress(buf.Bytes())
```

## 使用

```go
package main

import (
	"bytes"

	compress "github.com/vicanso/elton-compress/v2"
	"github.com/vicanso/elton/v2"
	"github.com/vicanso/elton/v2/middleware"
)

func main() {
	e := elton.New()
	// 需要注意添加的顺序，压缩是按添加的顺序选择适合的压缩方式
	// 此处只是示例所有的压缩器，正常使用时，按需使用1，2个压缩方式则可
	config := middleware.NewCompressConfig(
		// gzip、br、zstd 由 elton 内置
		&middleware.BrCompressor{
			MinLength: 1024,
		},
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

	e.GET("/", func(c *elton.Context) error {
		b := new(bytes.Buffer)
		for i := 0; i < 1000; i++ {
			b.WriteString("Hello, World!\n")
		}
		c.SetContentTypeByExt(".txt")
		c.BodyBuffer = b
		return nil
	})

	err := e.ListenAndServe(":3000")
	if err != nil {
		panic(err)
	}
}
```

```bash
curl -H 'Accept-Encoding: lz4' -v 'http://127.0.0.1:3000'
```

## Example

[middleware example](./example/main.go)

## 压缩率

jquery-3.4.1.min.js（88145 bytes）：

| 压缩方式 | 级别 | 大小 | 压缩率 |
|:-:|:-:|:-:|:-:|
| br | 6 | 29897 | 2.948 |
| gzip | 6 | 30827 | 2.859 |
| zstd | 3 | 32160 | 2.741 |
| lz4 | 1 | 46781 | 1.884 |
| lz4 | 4 | 39436 | 2.235 |
| snappy | 1 | 48127 | 1.832 |
| snappy | 2 | 42275 | 2.085 |
| snappy | 3 | 40026 | 2.202 |
| s2 | 1 | 48109 | 1.832 |
| s2 | 2 | 42002 | 2.099 |
| s2 | 3 | 39909 | 2.209 |

本仓库的压缩速度基准（`make bench`，不依赖网络）：

```bash
make bench
```
