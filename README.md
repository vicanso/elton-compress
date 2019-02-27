# cod-compress

[![Build Status](https://img.shields.io/travis/vicanso/cod-compress.svg?label=linux+build)](https://travis-ci.org/vicanso/cod-compress)

Compress middleware for cod, it support gzip compress function by default. For better performance, you can add more compress function such as brotli.

```go
package main

import (
	"bytes"

	"github.com/vicanso/cod"
	compress "github.com/vicanso/cod-compress"
)

func main() {
	d := cod.New()

	d.Use(compress.NewDefault())

	d.GET("/", func(c *cod.Context) (err error) {
		b := new(bytes.Buffer)
		for i := 0; i < 1000; i++ {
			b.WriteString("hello world\n")
		}
		c.SetHeader(cod.HeaderContentType, "text/plain; charset=utf-8")
		c.BodyBuffer = b
		return
	})

	d.ListenAndServe(":7001")
}
```