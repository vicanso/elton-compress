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
	"github.com/vicanso/elton/v2"
	"github.com/vicanso/elton/v2/middleware"
)

const (
	// S2Encoding s2 encoding
	S2Encoding = "s2"

	// S2LevelFastest 最快压缩级别（默认）
	S2LevelFastest = 1
	// S2LevelBetter 平衡压缩级别
	S2LevelBetter = 2
	// S2LevelBest 最高压缩级别
	S2LevelBest = 3
)

type (
	// S2Compressor s2 compress，输出为 s2 block 格式。
	// s2 是 snappy 的超集：压缩率与速度均优于 snappy，
	// 但产生的数据只能由 s2 解码器解压，snappy 解码器无法读取
	S2Compressor struct {
		// Level 压缩级别，1(最快) - 3(最高)，默认为 1
		Level int
		// MinLength 最小压缩长度，默认为 middleware.DefaultCompressMinLength
		MinLength int
	}
)

// NewS2Compressor returns a new s2 compressor
func NewS2Compressor() *S2Compressor {
	return &S2Compressor{}
}

func (s *S2Compressor) getLevel() int {
	return normalizeLevel(s.Level, S2LevelFastest, S2LevelBest)
}

func (s *S2Compressor) getMinLength() int {
	if s.MinLength == 0 {
		return middleware.DefaultCompressMinLength
	}
	return s.MinLength
}

// Accept check accept encoding
func (s *S2Compressor) Accept(c *elton.Context, bodySize int) (acceptable bool, encoding string) {
	// 如果数据少于最低压缩长度，则不压缩（reader 的 bodySize 为 -1，因此需要判断 >= 0）
	if bodySize >= 0 && bodySize < s.getMinLength() {
		return
	}
	return middleware.AcceptEncoding(c, S2Encoding)
}

// S2Compress compress data by s2 block format
func S2Compress(buf []byte, level int) (*bytes.Buffer, error) {
	var data []byte
	switch normalizeLevel(level, S2LevelFastest, S2LevelBest) {
	case S2LevelBetter:
		data = s2.EncodeBetter(nil, buf)
	case S2LevelBest:
		data = s2.EncodeBest(nil, buf)
	default:
		data = s2.Encode(nil, buf)
	}
	return bytes.NewBuffer(data), nil
}

// S2Decompress decompress data of s2 block format
func S2Decompress(buf []byte) (*bytes.Buffer, error) {
	data, err := s2.Decode(nil, buf)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

// Compress s2 compress
func (s *S2Compressor) Compress(buf []byte, levels ...int) (*bytes.Buffer, error) {
	level := s.getLevel()
	if len(levels) != 0 && levels[0] != middleware.IgnoreCompression {
		level = levels[0]
	}
	return S2Compress(buf, level)
}

// Pipe s2 pipe
func (s *S2Compressor) Pipe(c *elton.Context) error {
	// block 格式无法流式输出，一次读取所有数据再压缩
	r, closeReader := pipeReader(c)
	defer closeReader()

	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	buffer, err := S2Compress(buf, s.getLevel())
	if err != nil {
		return err
	}
	_, err = c.Response.Write(buffer.Bytes())
	return err
}
