// MIT License

// Copyright (c) 2020 Tree Xie

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package compress

import (
	"bytes"
	"io"

	"github.com/klauspost/compress/zstd"
	"github.com/vicanso/elton"
	"github.com/vicanso/elton/middleware"
)

const (
	// https://en.wikipedia.org/wiki/Zstd
	// In 2018 the algorithm was published as RFC 8478, which also defines an associated media type "application/zstd", filename extension "zst", and HTTP content encoding "zstd".[15]
	// ZstdEncoding zstd encoding
	ZstdEncoding = "zstd"
)

type (
	// ZstdCompressor zstd compress
	ZstdCompressor struct {
		Level     int
		MinLength int
	}
)

func (z *ZstdCompressor) getMinLength() int {
	if z.MinLength == 0 {
		return middleware.DefaultCompressMinLength
	}
	return z.MinLength
}

func (z *ZstdCompressor) getLevel() zstd.EncoderLevel {
	level := z.Level
	l := zstd.EncoderLevel(level)
	if l < zstd.SpeedFastest || l > zstd.SpeedBestCompression {
		return zstd.SpeedDefault
	}
	return l
}

// Accept check accept encoding
func (z *ZstdCompressor) Accept(c *elton.Context, bodySize int) (acceptable bool, encoding string) {
	// 如果数据少于最低压缩长度，则不压缩
	if bodySize >= 0 && bodySize < z.getMinLength() {
		return
	}
	return middleware.AcceptEncoding(c, ZstdEncoding)
}

// Compress zstd compress
func (z *ZstdCompressor) Compress(buf []byte, levels ...int) (*bytes.Buffer, error) {
	level := z.getLevel()
	if len(levels) != 0 && levels[0] != middleware.IgnoreCompression {
		level = zstd.EncoderLevel(levels[0])
	}
	encoder, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(level))
	if err != nil {
		return nil, err
	}
	data := encoder.EncodeAll(buf, make([]byte, 0, len(buf)))
	return bytes.NewBuffer(data), nil
}

// Pipe zstd pike
func (z *ZstdCompressor) Pipe(c *elton.Context) (err error) {
	r := c.Body.(io.Reader)
	closer, ok := c.Body.(io.Closer)
	if ok {
		defer closer.Close()
	}
	w, err := zstd.NewWriter(c.Response, zstd.WithEncoderLevel(z.getLevel()))
	if err != nil {
		return err
	}

	defer w.Close()
	_, err = io.Copy(w, r)
	return
}
