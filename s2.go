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

	"github.com/klauspost/compress/s2"
	"github.com/vicanso/elton"
	"github.com/vicanso/elton/middleware"
)

const (
	// S2Encoding s2 encoding
	S2Encoding = "s2"
)

type (
	// S2Compressor s2 compressor
	S2Compressor struct {
		Level     int
		MinLength int
	}
)

func (s *S2Compressor) getMinLength() int {
	if s.MinLength == 0 {
		return middleware.DefaultCompressMinLength
	}
	return s.MinLength
}

// Accept check accept encoding
func (s *S2Compressor) Accept(c *elton.Context, bodySize int) (acceptable bool, encoding string) {
	// 如果数据少于最低压缩长度，则不压缩
	if bodySize >= 0 && bodySize < s.getMinLength() {
		return
	}
	return middleware.AcceptEncoding(c, S2Encoding)
}

func s2IsBetterCompress(level int) bool {
	if level == 0 || level > 2 {
		return true
	}
	return false
}

// Compress s2 compress
func (s *S2Compressor) Compress(buf []byte) (*bytes.Buffer, error) {
	var dst []byte
	fn := s2.Encode
	if s2IsBetterCompress(s.Level) {
		fn = s2.EncodeBetter
	}
	data := fn(dst, buf)
	return bytes.NewBuffer(data), nil
}

// Pipe s2 pipe
func (s *S2Compressor) Pipe(c *elton.Context) (err error) {
	r := c.Body.(io.Reader)
	closer, ok := c.Body.(io.Closer)
	if ok {
		defer closer.Close()
	}
	var w *s2.Writer
	if s2IsBetterCompress(s.Level) {
		w = s2.NewWriter(c.Response, s2.WriterBetterCompression())
	} else {
		w = s2.NewWriter(c.Response)
	}
	defer w.Close()
	_, err = io.Copy(w, r)
	return
}
