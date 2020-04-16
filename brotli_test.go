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
	"io/ioutil"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")

// randomString get random string
func randomString(n int) string {
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func decodeBrotli(buf []byte) ([]byte, error) {
	r := brotli.NewReader(bytes.NewBuffer(buf))
	return ioutil.ReadAll(r)
}

func TestBrotliCompress(t *testing.T) {
	assert := assert.New(t)
	br := new(BrCompressor)
	originalData := randomString(1024)
	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	c := elton.NewContext(nil, req)
	acceptable, encoding := br.Accept(c, 0)
	assert.False(acceptable)
	assert.Empty(encoding)
	acceptable, encoding = br.Accept(c, len(originalData))
	assert.True(acceptable)
	assert.Equal(BrEncoding, encoding)
	buf, err := br.Compress([]byte(originalData))
	assert.Nil(err)
	originalBuf, _ := decodeBrotli(buf.Bytes())
	assert.Equal(originalData, string(originalBuf))
}

func TestBrotliPipe(t *testing.T) {
	assert := assert.New(t)
	resp := httptest.NewRecorder()
	originalData := randomString(1024)
	c := elton.NewContext(resp, nil)

	c.Body = bytes.NewReader([]byte(originalData))

	br := new(BrCompressor)
	err := br.Pipe(c)
	assert.Nil(err)
	buf, _ := decodeBrotli(resp.Body.Bytes())
	assert.Equal(originalData, string(buf))
}
