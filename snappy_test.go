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
	"net/http/httptest"
	"testing"

	"github.com/klauspost/compress/snappy"
	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton/v2"
	"github.com/vicanso/elton/v2/middleware"
)

func TestSnappyAccept(t *testing.T) {
	assert := assert.New(t)
	s := NewSnappyCompressor()

	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", SnappyEncoding)
	c := elton.NewContext(nil, req)

	acceptable, encoding := s.Accept(c, 0)
	assert.False(acceptable)
	assert.Empty(encoding)

	acceptable, encoding = s.Accept(c, len(compressTestData))
	assert.True(acceptable)
	assert.Equal(SnappyEncoding, encoding)

	// 自定义最小压缩长度
	acceptable, _ = (&SnappyCompressor{MinLength: 10 * 1024}).Accept(c, len(compressTestData))
	assert.False(acceptable)
}

func TestSnappyCompress(t *testing.T) {
	assert := assert.New(t)

	buf, err := NewSnappyCompressor().Compress(compressTestData)
	assert.Nil(err)
	assert.NotEmpty(buf)

	// snappy block 自带原始长度前缀，标准 snappy 解码器可直接解压
	originalBuf, err := snappy.Decode(nil, buf.Bytes())
	assert.Nil(err)
	assert.Equal(compressTestData, originalBuf)

	dst, err := SnappyDecompress(buf.Bytes())
	assert.Nil(err)
	assert.Equal(compressTestData, dst.Bytes())
}

func TestSnappyCompressLevel(t *testing.T) {
	assert := assert.New(t)

	// 零值 Level 的默认档位为 Better，与 v1 的 snappy.Encode 行为一致
	defaultBuf, err := NewSnappyCompressor().Compress(compressTestData)
	assert.Nil(err)
	betterBuf, err := SnappyCompress(compressTestData, SnappyLevelBetter)
	assert.Nil(err)
	assert.Equal(betterBuf.Bytes(), defaultBuf.Bytes())

	fastest, err := SnappyCompress(compressTestData, SnappyLevelFastest)
	assert.Nil(err)
	best, err := SnappyCompress(compressTestData, SnappyLevelBest)
	assert.Nil(err)
	assert.Less(best.Len(), fastest.Len())

	// 越界级别收敛至上限
	buf, err := SnappyCompress(compressTestData, 100)
	assert.Nil(err)
	assert.Equal(best.Bytes(), buf.Bytes())

	// IgnoreCompression 时使用压缩器自身的级别
	buf, err = (&SnappyCompressor{Level: SnappyLevelBest}).Compress(compressTestData, middleware.IgnoreCompression)
	assert.Nil(err)
	assert.Equal(best.Bytes(), buf.Bytes())

	// DynamicLevel 传入的级别覆盖压缩器自身的级别
	buf, err = (&SnappyCompressor{Level: SnappyLevelFastest}).Compress(compressTestData, SnappyLevelBest)
	assert.Nil(err)
	assert.Equal(best.Bytes(), buf.Bytes())
}

func TestSnappyPipe(t *testing.T) {
	assert := assert.New(t)
	resp := httptest.NewRecorder()
	c := elton.NewContext(resp, nil)
	body := &closeAwareReader{Reader: bytes.NewReader(compressTestData)}
	c.Body = body

	err := NewSnappyCompressor().Pipe(c)
	assert.Nil(err)
	assert.True(body.closed)
	assert.NotEmpty(resp.Body.Bytes())

	dst, err := SnappyDecompress(resp.Body.Bytes())
	assert.Nil(err)
	assert.Equal(compressTestData, dst.Bytes())
}

func TestSnappyPipeError(t *testing.T) {
	assert := assert.New(t)
	c := elton.NewContext(httptest.NewRecorder(), nil)
	c.Body = &closeAwareReader{Reader: &errReader{}}

	err := NewSnappyCompressor().Pipe(c)
	assert.ErrorIs(err, errRead)
}

func TestSnappyDecompress(t *testing.T) {
	assert := assert.New(t)

	_, err := SnappyDecompress([]byte("not a snappy block"))
	assert.NotNil(err)
}
