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
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/snappy"
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

func TestSnappyCompress(t *testing.T) {
	assert := assert.New(t)
	originalData := randomString(1024)
	s := new(SnappyCompressor)

	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", SnappyEncoding)
	c := elton.NewContext(nil, req)
	acceptable, encoding := s.Accept(c, 0)
	assert.False(acceptable)
	assert.Empty(encoding)
	acceptable, encoding = s.Accept(c, len(originalData))
	assert.True(acceptable)
	assert.Equal(encoding, SnappyEncoding)

	buf, err := s.Compress([]byte(originalData))
	assert.Nil(err)
	assert.NotEmpty(buf)
	var originalBuf []byte
	originalBuf, err = snappy.Decode(originalBuf, buf.Bytes())
	assert.Nil(err)
	assert.Equal(originalData, string(originalBuf))
}

func TestSnappyPipe(t *testing.T) {
	assert := assert.New(t)
	resp := httptest.NewRecorder()
	originalData := randomString(1024)
	c := elton.NewContext(resp, nil)

	c.Body = bytes.NewReader([]byte(originalData))

	s := new(SnappyCompressor)
	err := s.Pipe(c)
	assert.Nil(err)
	assert.NotEmpty(resp.Body.Bytes())

	var originalBuf []byte
	originalBuf, err = snappy.Decode(originalBuf, resp.Body.Bytes())
	assert.Nil(err)
	assert.Equal(originalData, string(originalBuf))
}
