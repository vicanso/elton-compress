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
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton/v2"
	"github.com/vicanso/elton/v2/middleware"
)

var errRead = errors.New("read fail")

// errReader 是一个总是读取失败的 reader
type errReader struct{}

func (r *errReader) Read([]byte) (int, error) {
	return 0, errRead
}

func TestNormalizeLevel(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		level  int
		expect int
	}{
		// 非正值（含零值 Level 与 middleware.IgnoreCompression）回退到默认级别
		{level: 0, expect: 2},
		{level: -1, expect: 2},
		{level: middleware.IgnoreCompression, expect: 2},
		{level: 1, expect: 1},
		{level: 3, expect: 3},
		// 超出上限收敛至上限
		{level: 4, expect: 3},
		{level: 1000, expect: 3},
	}
	for _, tt := range tests {
		assert.Equal(tt.expect, normalizeLevel(tt.level, 2, 3))
	}
}

func TestPipeReader(t *testing.T) {
	assert := assert.New(t)

	// body 实现了 io.Closer，Pipe 结束后应被关闭
	c := elton.NewContext(nil, nil)
	body := &closeAwareReader{Reader: bytes.NewReader(compressTestData)}
	c.Body = body
	r, closeReader := pipeReader(c)
	assert.Equal(body, r)
	assert.False(body.closed)
	closeReader()
	assert.True(body.closed)

	// body 未实现 io.Closer 时不应 panic
	c = elton.NewContext(nil, nil)
	c.Body = bytes.NewReader(compressTestData)
	_, closeReader = pipeReader(c)
	assert.NotPanics(closeReader)
}

// TestCompressorInterface 确保三个压缩器都满足 elton 的 Compressor 接口，
// 且可以被 compress 中间件正常选中
func TestCompressorInterface(t *testing.T) {
	assert := assert.New(t)
	compressors := map[string]middleware.Compressor{
		SnappyEncoding: NewSnappyCompressor(),
		S2Encoding:     NewS2Compressor(),
		Lz4Encoding:    NewLz4Compressor(),
	}
	decompressors := map[string]func([]byte) (*bytes.Buffer, error){
		SnappyEncoding: SnappyDecompress,
		S2Encoding:     S2Decompress,
		Lz4Encoding:    Lz4Decompress,
	}

	for encoding, compressor := range compressors {
		fn := middleware.NewCompress(middleware.NewCompressConfig(compressor))

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set(elton.HeaderAcceptEncoding, encoding)
		c := elton.NewContext(httptest.NewRecorder(), req)
		c.Next = func() error {
			c.SetHeader(elton.HeaderContentType, "text/plain; charset=utf-8")
			c.BodyBuffer = bytes.NewBuffer(compressTestData)
			return nil
		}

		err := fn(c)
		assert.Nil(err, encoding)
		assert.Equal(encoding, c.GetHeader(elton.HeaderContentEncoding), encoding)
		assert.Less(c.BodyBuffer.Len(), len(compressTestData), encoding)

		data, err := decompressors[encoding](c.BodyBuffer.Bytes())
		assert.Nil(err, encoding)
		assert.Equal(compressTestData, data.Bytes(), encoding)
	}
}
