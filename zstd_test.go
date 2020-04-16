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

	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

func TestZstdCompress(t *testing.T) {
	assert := assert.New(t)
	originalData := randomString(1024)
	z := new(ZstdCompressor)

	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", ZstdEncoding)
	c := elton.NewContext(nil, req)
	acceptable, encoding := z.Accept(c, len(originalData))
	assert.True(acceptable)
	assert.Equal(encoding, ZstdEncoding)

	buf, err := z.Compress([]byte(originalData))
	assert.Nil(err)
	assert.NotEmpty(buf)

	decorder, err := zstd.NewReader(nil)
	assert.Nil(err)
	var dst []byte
	dst, err = decorder.DecodeAll(buf.Bytes(), dst)
	assert.Nil(err)
	assert.Equal([]byte(originalData), dst)
}

func TestZstdPipe(t *testing.T) {
	assert := assert.New(t)
	resp := httptest.NewRecorder()
	originalData := randomString(1024)
	c := elton.NewContext(resp, nil)
	c.Body = bytes.NewReader([]byte(originalData))

	z := new(ZstdCompressor)
	err := z.Pipe(c)
	assert.Nil(err)
	assert.NotEmpty(resp.Body.Bytes())

	decorder, err := zstd.NewReader(nil)
	assert.Nil(err)
	var dst []byte
	dst, err = decorder.DecodeAll(resp.Body.Bytes(), dst)
	assert.Nil(err)
	assert.Equal([]byte(originalData), dst)
}
