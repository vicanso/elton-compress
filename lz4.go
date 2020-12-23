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
	"errors"
	"io"
	"io/ioutil"

	"github.com/pierrec/lz4"
	"github.com/vicanso/elton"
	"github.com/vicanso/elton/middleware"
)

const (
	// Lz4Encoding lz4 encoding
	Lz4Encoding = "lz4"
)

type (
	// Lz4Compressor lz4 compress
	Lz4Compressor struct {
		Level     int
		MinLength int
	}
)

var ErrLz4IsNotCompressible = errors.New("Is not compressible for lz4")

func (l *Lz4Compressor) getMinLength() int {
	if l.MinLength == 0 {
		return middleware.DefaultCompressMinLength
	}
	return l.MinLength
}

// Accept check accept encoding
func (l *Lz4Compressor) Accept(c *elton.Context, bodySize int) (acceptable bool, encoding string) {
	// 如果数据少于最低压缩长度，则不压缩
	if bodySize >= 0 && bodySize < l.getMinLength() {
		return
	}
	return middleware.AcceptEncoding(c, Lz4Encoding)
}

// Compress lz4 compress
func (l *Lz4Compressor) Compress(buf []byte) (*bytes.Buffer, error) {
	dst := make([]byte, len(buf))
	n, err := lz4.CompressBlock(buf, dst, nil)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, ErrLz4IsNotCompressible
	}
	return bytes.NewBuffer(dst[:n]), nil
}

// Pipe lz4 pipe compress
func (l *Lz4Compressor) Pipe(c *elton.Context) (err error) {
	// 使用lz4时为了提升性能，还是使用compress block的方式
	// 一次读取所有数据
	r := c.Body.(io.Reader)
	closer, ok := c.Body.(io.Closer)
	if ok {
		defer closer.Close()
	}
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	buffer, err := l.Compress(buf)
	if err != nil {
		return
	}
	_, err = c.Response.Write(buffer.Bytes())
	return
}
