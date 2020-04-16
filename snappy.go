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

	"github.com/golang/snappy"
	"github.com/vicanso/elton"
	"github.com/vicanso/elton/middleware"
)

const (
	// SnappyEncoding snappy encoding
	SnappyEncoding = "snz"
)

type (
	// SnappyCompressor snappy compress
	SnappyCompressor struct {
		MinLength int
	}
)

func (s *SnappyCompressor) getMinLength() int {
	if s.MinLength == 0 {

		return middleware.DefaultCompressMinLength
	}
	return s.MinLength
}

// Accept check accept encoding
func (s *SnappyCompressor) Accept(c *elton.Context, bodySize int) (acceptable bool, encoding string) {
	// 如果数据少于最低压缩长度，则不压缩
	if bodySize >= 0 && bodySize < s.getMinLength() {
		return
	}
	return middleware.AcceptEncoding(c, SnappyEncoding)
}

// Compress snappy compress
func (s *SnappyCompressor) Compress(buf []byte) (*bytes.Buffer, error) {
	var dst []byte
	data := snappy.Encode(dst, buf)
	return bytes.NewBuffer(data), nil
}

// Pipe snappy pipe
func (s *SnappyCompressor) Pipe(c *elton.Context) (err error) {
	r := c.Body.(io.Reader)
	closer, ok := c.Body.(io.Closer)
	if ok {
		defer closer.Close()
	}
	w := snappy.NewBufferedWriter(c.Response)
	defer w.Close()
	_, err = io.Copy(w, r)
	return
}
