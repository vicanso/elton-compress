# cod-compress

Compress middleware for cod.

```go
package main

import (
	"bytes"

	"github.com/vicanso/cod"
	compress "github.com/vicanso/cod-compress"
)

func main() {
	d := cod.New()
	d.Keys = []string{
		"cuttlefish",
	}
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