# elton-compress

[![Build Status](https://img.shields.io/travis/vicanso/elton-compress.svg?label=linux+build)](https://travis-ci.org/vicanso/elton-compress)

Compress middleware for elton, it support gzip and br compress function by default. 

- `BrCompressor` brotli compress is better for http, most modern browser support it.
- `SnappyCompressor` snappy compress is fast, but not aim for maximum compression. It's useful for Intranet.


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

## brotli

I use [Pure Go Brotli](https://github.com/andybalholm/brotli) instead of cbrotli, it is a little slower. 

Test for encode html(113K) to br:

```bash
BenchmarkCBrotli-8    	     200	   9555795 ns/op
BenchmarkGoBrotli-8   	     100	  10703582 ns/op
```