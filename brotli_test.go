package compress

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

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
	acceptable, encoding := br.Accept(c)
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
