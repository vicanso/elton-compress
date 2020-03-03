package compress

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

func TestGzipCompress(t *testing.T) {
	assert := assert.New(t)
	originalData := randomString(1024)
	g := new(GzipCompressor)
	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	c := elton.NewContext(nil, req)
	acceptable, encoding := g.Accept(c, 0)
	assert.False(acceptable)
	assert.Empty(encoding)
	acceptable, encoding = g.Accept(c, len(originalData))
	assert.True(acceptable)
	assert.Equal(GzipEncoding, encoding)
	buf, err := g.Compress([]byte(originalData))
	assert.Nil(err)
	r, err := gzip.NewReader(bytes.NewReader(buf.Bytes()))
	assert.Nil(err)
	defer r.Close()
	originlBuf, _ := ioutil.ReadAll(r)
	assert.Equal(originalData, string(originlBuf))
}

func TestGzipPipe(t *testing.T) {
	assert := assert.New(t)
	resp := httptest.NewRecorder()
	originalData := randomString(1024)
	c := elton.NewContext(resp, nil)

	c.Body = bytes.NewReader([]byte(originalData))

	g := new(GzipCompressor)
	err := g.Pipe(c)
	assert.Nil(err)
	r, err := gzip.NewReader(resp.Body)
	assert.Nil(err)
	defer r.Close()
	buf, _ := ioutil.ReadAll(r)
	assert.Equal(originalData, string(buf))
}
