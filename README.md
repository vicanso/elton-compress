# elton-compress

[![Build Status](https://img.shields.io/travis/vicanso/elton-compress.svg?label=linux+build)](https://travis-ci.org/vicanso/elton-compress)

Compress middleware for elton, it support gzip and br compress function by default. 

- `BrCompressor` Brotli compression algorithm is better for http, most modern browser support it. Compress level is 1-11，default 0(6).
- `GzipComperssor` Gzip compression algorithm is most commonly used for http, all modern browser support it. Compress level is 1-9, default 0(6).
- `SnappyCompressor` Snappy compression algorithm is fast, but not aim for maximum compression. It's useful for Intranet. Not support compress level.
- `ZstdCompressor` Zstandard is a real-time compression algorithm, providing high compression ratios. Compress level is 1-2, default 0(2).
- `S2Compressor` S2 is a high performance replacement for Snappy. Compress level is 1-2, default 0(2).
- `Lz4Compressor` LZ4 is lossless compression algorithm, providing compression speed > 500 MB/s per core, scalable with multi-cores CPU. Compress level higher is better, use 0 for fastest compression.

```go
package main

import (
	"bytes"

	"github.com/vicanso/elton"
	compress "github.com/vicanso/elton-compress"
)

func main() {
	d := elton.New()

	d.Use(compress.NewDefault())

	d.GET("/", func(c *elton.Context) (err error) {
		b := new(bytes.Buffer)
		for i := 0; i < 1000; i++ {
			b.WriteString("hello world\n")
		}
		c.SetHeader(elton.HeaderContentType, "text/plain; charset=utf-8")
		c.BodyBuffer = b
		return
	})

	d.ListenAndServe(":7001")
}
```

## Example

[middleware example](./example/main.go)

jquery-3.4.1.min.js(88145 bytes)

compression | size | ratio | level
:-:|:-:|:-:|:-:
gzip | 30827 | 2.859 | 6 
br | 29897 | 2.948 | 6
snappy | 47709 | 1.847 | - 
zstd | 32816 | 2.686 | 2
s2 | 42750 | 2.061 | 2
lz4 | 39434 | 2.235 | 6

## brotli

I use [Pure Go Brotli](https://github.com/andybalholm/brotli) instead of cbrotli, it is a little bit slow. 

Test for encode html(113K) to br:

```bash
BenchmarkCBrotli-8    	     200	   9555795 ns/op
BenchmarkGoBrotli-8   	     100	  10703582 ns/op
```