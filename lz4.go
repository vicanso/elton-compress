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

	"github.com/pierrec/lz4/v4"
	"github.com/vicanso/elton/v2"
	"github.com/vicanso/elton/v2/middleware"
)

const (
	// Lz4Encoding lz4 encoding
	Lz4Encoding = "lz4"

	// Lz4LevelFastest 最快压缩级别（对应 lz4.Fast）
	Lz4LevelFastest = 1
	// Lz4LevelBest 最高压缩级别（对应 lz4.Level9）
	Lz4LevelBest = 10
)

// lz4Levels 把 [Lz4LevelFastest, Lz4LevelBest] 映射到 lz4 的级别常量。
// lz4.Level1 起为 1<<8 这样的位标记而非连续整数，无法直接用整数换算。
var lz4Levels = [...]lz4.CompressionLevel{
	lz4.Fast,
	lz4.Level1,
	lz4.Level2,
	lz4.Level3,
	lz4.Level4,
	lz4.Level5,
	lz4.Level6,
	lz4.Level7,
	lz4.Level8,
	lz4.Level9,
}

type (
	// Lz4Compressor lz4 compress
	Lz4Compressor struct {
		// Level 压缩级别，1(最快) - 10(最高)，默认为 1
		Level int
		// MinLength 最小压缩长度，默认为 middleware.DefaultCompressMinLength
		MinLength int
	}
)

// NewLz4Compressor returns a new lz4 compressor
func NewLz4Compressor() *Lz4Compressor {
	return &Lz4Compressor{}
}

func (l *Lz4Compressor) getLevel() int {
	return normalizeLevel(l.Level, Lz4LevelFastest, Lz4LevelBest)
}

func (l *Lz4Compressor) getMinLength() int {
	if l.MinLength == 0 {
		return middleware.DefaultCompressMinLength
	}
	return l.MinLength
}

// Accept check accept encoding
func (l *Lz4Compressor) Accept(c *elton.Context, bodySize int) (acceptable bool, encoding string) {
	// 如果数据少于最低压缩长度，则不压缩（reader 的 bodySize 为 -1，因此需要判断 >= 0）
	if bodySize >= 0 && bodySize < l.getMinLength() {
		return
	}
	return middleware.AcceptEncoding(c, Lz4Encoding)
}

// Lz4Compress compress data by lz4 frame format
func Lz4Compress(buf []byte, level int) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	w, err := newLz4Writer(buffer, level)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(buf)
	if err != nil {
		return nil, err
	}
	// 直接调用 close 触发数据的 flush
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

// Lz4Decompress decompress data of lz4 frame format
func Lz4Decompress(buf []byte) (*bytes.Buffer, error) {
	data, err := io.ReadAll(lz4.NewReader(bytes.NewReader(buf)))
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

// Compress lz4 compress
func (l *Lz4Compressor) Compress(buf []byte, levels ...int) (*bytes.Buffer, error) {
	level := l.getLevel()
	if len(levels) != 0 && levels[0] != middleware.IgnoreCompression {
		level = levels[0]
	}
	return Lz4Compress(buf, level)
}

// Pipe lz4 pipe compress
func (l *Lz4Compressor) Pipe(c *elton.Context) error {
	r, closeReader := pipeReader(c)
	defer closeReader()

	w, err := newLz4Writer(c.Response, l.getLevel())
	if err != nil {
		return err
	}
	_, err = io.Copy(w, r)
	if err != nil {
		_ = w.Close()
		return err
	}
	// close 会写入 frame 的结束标记，其错误不可忽略
	return w.Close()
}

func newLz4Writer(w io.Writer, level int) (*lz4.Writer, error) {
	writer := lz4.NewWriter(w)
	compressionLevel := lz4Levels[normalizeLevel(level, Lz4LevelFastest, Lz4LevelBest)-1]
	err := writer.Apply(lz4.CompressionLevelOption(compressionLevel))
	if err != nil {
		return nil, err
	}
	return writer, nil
}
