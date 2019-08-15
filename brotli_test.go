// +build brotli

package compress

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/google/brotli/go/cbrotli"
	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

func TestBrotliCompress(t *testing.T) {
	assert := assert.New(t)
	br := new(BrCompressor)
	originalData := randomString(1024)
	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	c := elton.NewContext(nil, req)
	acceptable, encoding := br.Accept(c)
	assert.True(acceptable)
	assert.Equal(brEncoding, encoding)
	buf, err := br.Compress([]byte(originalData), 0)
	assert.Nil(err)
	originalBuf, _ := cbrotli.Decode(buf)
	assert.Equal(originalData, string(originalBuf))
}

func TestBrotliPipe(t *testing.T) {
	assert := assert.New(t)
	resp := httptest.NewRecorder()
	originalData := randomString(1024)
	c := elton.NewContext(resp, nil)

	c.Body = bytes.NewReader([]byte(originalData))

	br := new(BrCompressor)
	err := br.Pipe(c, 0)
	assert.Nil(err)
	buf, _ := cbrotli.Decode(resp.Body.Bytes())
	assert.Equal(originalData, string(buf))
}
