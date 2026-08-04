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

	"github.com/klauspost/compress/s2"
	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton/v2"
	"github.com/vicanso/elton/v2/middleware"
)

func TestS2Accept(t *testing.T) {
	assert := assert.New(t)
	s := NewS2Compressor()

	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", S2Encoding)
	c := elton.NewContext(nil, req)

	acceptable, encoding := s.Accept(c, 0)
	assert.False(acceptable)
	assert.Empty(encoding)

	acceptable, encoding = s.Accept(c, len(compressTestData))
	assert.True(acceptable)
	assert.Equal(S2Encoding, encoding)

	// 自定义最小压缩长度
	acceptable, _ = (&S2Compressor{MinLength: 10 * 1024}).Accept(c, len(compressTestData))
	assert.False(acceptable)
}

func TestS2Compress(t *testing.T) {
	assert := assert.New(t)

	buf, err := NewS2Compressor().Compress(compressTestData)
	assert.Nil(err)
	assert.NotEmpty(buf)

	originalBuf, err := s2.Decode(nil, buf.Bytes())
	assert.Nil(err)
	assert.Equal(compressTestData, originalBuf)

	dst, err := S2Decompress(buf.Bytes())
	assert.Nil(err)
	assert.Equal(compressTestData, dst.Bytes())
}

func TestS2CompressLevel(t *testing.T) {
	assert := assert.New(t)

	// 零值 Level 的默认档位为 Fastest
	defaultBuf, err := NewS2Compressor().Compress(compressTestData)
	assert.Nil(err)
	fastest, err := S2Compress(compressTestData, S2LevelFastest)
	assert.Nil(err)
	assert.Equal(fastest.Bytes(), defaultBuf.Bytes())

	better, err := S2Compress(compressTestData, S2LevelBetter)
	assert.Nil(err)
	assert.Equal(s2.EncodeBetter(nil, compressTestData), better.Bytes())

	best, err := S2Compress(compressTestData, S2LevelBest)
	assert.Nil(err)
	assert.Less(best.Len(), fastest.Len())

	// 越界级别收敛至上限
	buf, err := S2Compress(compressTestData, 100)
	assert.Nil(err)
	assert.Equal(best.Bytes(), buf.Bytes())

	// IgnoreCompression 时使用压缩器自身的级别
	buf, err = (&S2Compressor{Level: S2LevelBest}).Compress(compressTestData, middleware.IgnoreCompression)
	assert.Nil(err)
	assert.Equal(best.Bytes(), buf.Bytes())

	// DynamicLevel 传入的级别覆盖压缩器自身的级别
	buf, err = (&S2Compressor{Level: S2LevelFastest}).Compress(compressTestData, S2LevelBest)
	assert.Nil(err)
	assert.Equal(best.Bytes(), buf.Bytes())
}

func TestS2Pipe(t *testing.T) {
	assert := assert.New(t)
	resp := httptest.NewRecorder()
	c := elton.NewContext(resp, nil)
	body := &closeAwareReader{Reader: bytes.NewReader(compressTestData)}
	c.Body = body

	err := NewS2Compressor().Pipe(c)
	assert.Nil(err)
	assert.True(body.closed)
	assert.NotEmpty(resp.Body.Bytes())

	dst, err := S2Decompress(resp.Body.Bytes())
	assert.Nil(err)
	assert.Equal(compressTestData, dst.Bytes())
}

func TestS2PipeError(t *testing.T) {
	assert := assert.New(t)
	c := elton.NewContext(httptest.NewRecorder(), nil)
	c.Body = &closeAwareReader{Reader: &errReader{}}

	err := NewS2Compressor().Pipe(c)
	assert.ErrorIs(err, errRead)
}

func TestS2Decompress(t *testing.T) {
	assert := assert.New(t)

	_, err := S2Decompress([]byte("not a s2 block"))
	assert.NotNil(err)
}
