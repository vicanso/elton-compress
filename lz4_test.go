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
	"net/http/httptest"
	"testing"

	"github.com/pierrec/lz4/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton/v2"
	"github.com/vicanso/elton/v2/middleware"
)

var compressTestData = []byte(`Brotli is a data format specification[2] for data streams compressed with a specific combination of the general-purpose LZ77 lossless compression algorithm, Huffman coding and 2nd order context modelling. Brotli is a compression algorithm developed by Google and works best for text compression.

Google employees Jyrki Alakuijala and Zoltán Szabadka initially developed Brotli to decrease the size of transmissions of WOFF2 web fonts, and in that context Brotli was a continuation of the development of zopfli, which is a zlib-compatible implementation of the standard gzip and deflate specifications. Brotli allows a denser packing than gzip and deflate because of several algorithmic and format level improvements: the use of context models for literals and copy distances, describing copy distances through past distances, use of move-to-front queue in entropy code selection, joint-entropy coding of literal and copy lengths, the use of graph algorithms in block splitting, and a larger backward reference window are example improvements. The Brotli specification was generalized in September 2015 for HTTP stream compression (content-encoding type 'br'). This generalized iteration also improved the compression ratio by using a pre-defined dictionary of frequently used words and phrases.`)

func TestLz4Accept(t *testing.T) {
	assert := assert.New(t)
	l := NewLz4Compressor()

	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", Lz4Encoding)
	c := elton.NewContext(nil, req)

	// 少于最小压缩长度
	acceptable, encoding := l.Accept(c, 0)
	assert.False(acceptable)
	assert.Empty(encoding)

	acceptable, encoding = l.Accept(c, len(compressTestData))
	assert.True(acceptable)
	assert.Equal(Lz4Encoding, encoding)

	// reader body 的长度未知，仍然可压缩
	acceptable, encoding = l.Accept(c, -1)
	assert.True(acceptable)
	assert.Equal(Lz4Encoding, encoding)

	// 自定义最小压缩长度
	acceptable, _ = (&Lz4Compressor{MinLength: 10 * 1024}).Accept(c, len(compressTestData))
	assert.False(acceptable)

	// 客户端不支持 lz4
	req = httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", "gzip, br")
	acceptable, encoding = l.Accept(elton.NewContext(nil, req), len(compressTestData))
	assert.False(acceptable)
	assert.Empty(encoding)
}

func TestLz4Compress(t *testing.T) {
	assert := assert.New(t)
	l := NewLz4Compressor()

	buf, err := l.Compress(compressTestData)
	assert.Nil(err)
	assert.NotEmpty(buf)

	// 输出为标准 lz4 frame，无需外部提供原始长度即可解压
	valid, err := lz4.ValidFrameHeader(buf.Bytes())
	assert.Nil(err)
	assert.True(valid)

	dst, err := Lz4Decompress(buf.Bytes())
	assert.Nil(err)
	assert.Equal(compressTestData, dst.Bytes())
}

func TestLz4CompressLevel(t *testing.T) {
	assert := assert.New(t)

	// Level 字段生效：最高级别应压得比最快级别更小
	fastest, err := (&Lz4Compressor{Level: Lz4LevelFastest}).Compress(compressTestData)
	assert.Nil(err)
	best, err := (&Lz4Compressor{Level: Lz4LevelBest}).Compress(compressTestData)
	assert.Nil(err)
	assert.Less(best.Len(), fastest.Len())

	// 超出范围的级别收敛至边界，而不是报错
	tests := []struct {
		level  int
		expect int
	}{
		{level: 0, expect: Lz4LevelFastest},
		{level: middleware.IgnoreCompression, expect: Lz4LevelFastest},
		{level: 100, expect: Lz4LevelBest},
	}
	for _, tt := range tests {
		buf, err := Lz4Compress(compressTestData, tt.level)
		assert.Nil(err)
		expected, err := Lz4Compress(compressTestData, tt.expect)
		assert.Nil(err)
		assert.Equal(expected.Bytes(), buf.Bytes())
	}

	// 中间件传入的 IgnoreCompression 使用压缩器自身的级别
	buf, err := (&Lz4Compressor{Level: Lz4LevelBest}).Compress(compressTestData, middleware.IgnoreCompression)
	assert.Nil(err)
	assert.Equal(best.Bytes(), buf.Bytes())

	// 中间件 DynamicLevel 传入的级别覆盖压缩器自身的级别
	buf, err = (&Lz4Compressor{Level: Lz4LevelFastest}).Compress(compressTestData, Lz4LevelBest)
	assert.Nil(err)
	assert.Equal(best.Bytes(), buf.Bytes())
}

func TestLz4Pipe(t *testing.T) {
	assert := assert.New(t)
	resp := httptest.NewRecorder()
	c := elton.NewContext(resp, nil)
	body := &closeAwareReader{Reader: bytes.NewReader(compressTestData)}
	c.Body = body

	err := NewLz4Compressor().Pipe(c)
	assert.Nil(err)
	assert.True(body.closed)
	assert.NotEmpty(resp.Body.Bytes())

	dst, err := Lz4Decompress(resp.Body.Bytes())
	assert.Nil(err)
	assert.Equal(compressTestData, dst.Bytes())
}

func TestLz4PipeError(t *testing.T) {
	assert := assert.New(t)
	c := elton.NewContext(httptest.NewRecorder(), nil)
	c.Body = &closeAwareReader{Reader: &errReader{}}

	err := NewLz4Compressor().Pipe(c)
	assert.ErrorIs(err, errRead)
}

func TestLz4Decompress(t *testing.T) {
	assert := assert.New(t)

	_, err := Lz4Decompress([]byte("not a lz4 frame"))
	assert.NotNil(err)
}

// closeAwareReader 用于校验 Pipe 是否关闭了实现 io.Closer 的 body
type closeAwareReader struct {
	io.Reader
	closed bool
}

func (r *closeAwareReader) Close() error {
	r.closed = true
	return nil
}
