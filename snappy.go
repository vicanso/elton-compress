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
	"github.com/klauspost/compress/snappy"
	"github.com/vicanso/elton/v2"
	"github.com/vicanso/elton/v2/middleware"
)

const (
	// SnappyEncoding snappy encoding
	SnappyEncoding = "snz"

	// SnappyLevelFastest 最快压缩级别
	SnappyLevelFastest = 1
	// SnappyLevelBetter 平衡压缩级别（默认）
	SnappyLevelBetter = 2
	// SnappyLevelBest 最高压缩级别
	SnappyLevelBest = 3
)

type (
	// SnappyCompressor snappy compress，输出为 snappy block 格式，
	// block 自带原始长度前缀，客户端可直接解压
	SnappyCompressor struct {
		// Level 压缩级别，1(最快) - 3(最高)，默认为 2
		Level int
		// MinLength 最小压缩长度，默认为 middleware.DefaultCompressMinLength
		MinLength int
	}
)

// NewSnappyCompressor returns a new snappy compressor
func NewSnappyCompressor() *SnappyCompressor {
	return &SnappyCompressor{}
}

func (s *SnappyCompressor) getLevel() int {
	return normalizeLevel(s.Level, SnappyLevelBetter, SnappyLevelBest)
}

func (s *SnappyCompressor) getMinLength() int {
	if s.MinLength == 0 {
		return middleware.DefaultCompressMinLength
	}
	return s.MinLength
}

// Accept check accept encoding
func (s *SnappyCompressor) Accept(c *elton.Context, bodySize int) (acceptable bool, encoding string) {
	// 如果数据少于最低压缩长度，则不压缩（reader 的 bodySize 为 -1，因此需要判断 >= 0）
	if bodySize >= 0 && bodySize < s.getMinLength() {
		return
	}
	return middleware.AcceptEncoding(c, SnappyEncoding)
}

// SnappyCompress compress data by snappy block format
func SnappyCompress(buf []byte, level int) (*bytes.Buffer, error) {
	var data []byte
	switch normalizeLevel(level, SnappyLevelBetter, SnappyLevelBest) {
	case SnappyLevelFastest:
		data = s2.EncodeSnappy(nil, buf)
	case SnappyLevelBest:
		data = s2.EncodeSnappyBest(nil, buf)
	default:
		data = s2.EncodeSnappyBetter(nil, buf)
	}
	return bytes.NewBuffer(data), nil
}

// SnappyDecompress decompress data of snappy block format
func SnappyDecompress(buf []byte) (*bytes.Buffer, error) {
	data, err := snappy.Decode(nil, buf)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

// Compress snappy compress
func (s *SnappyCompressor) Compress(buf []byte, levels ...int) (*bytes.Buffer, error) {
	level := s.getLevel()
	if len(levels) != 0 && levels[0] != middleware.IgnoreCompression {
		level = levels[0]
	}
	return SnappyCompress(buf, level)
}

// Pipe snappy pipe
func (s *SnappyCompressor) Pipe(c *elton.Context) error {
	// block 格式无法流式输出，一次读取所有数据再压缩
	r, closeReader := pipeReader(c)
	defer closeReader()

	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	buffer, err := SnappyCompress(buf, s.getLevel())
	if err != nil {
		return err
	}
	_, err = c.Response.Write(buffer.Bytes())
	return err
}
