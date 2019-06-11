package compress

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/cod"
)

func TestDoGzip(t *testing.T) {
	assert := assert.New(t)
	buf := []byte("abcd")
	_, err := doGzip(buf, 0)
	assert.Nil(err)

	_, err = doGzip(buf, 100)
	assert.Nil(err)
}

func TestGzipCompress(t *testing.T) {
	assert := assert.New(t)
	originalData := randomString(1024)
	g := new(GzipCompressor)
	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	c := cod.NewContext(nil, req)
	acceptable, encoding := g.Accept(c)
	assert.True(acceptable)
	assert.Equal(encoding, gzipEncoding)
	buf, err := g.Compress([]byte(originalData), 0)
	assert.Nil(err)
	r, err := gzip.NewReader(bytes.NewBuffer(buf))
	assert.Nil(err)
	defer r.Close()
	originlBuf, _ := ioutil.ReadAll(r)
	assert.Equal(string(originlBuf), originalData)
}

// doGunzip gunzip
func doGunzip(buf []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

func TestGzipPipe(t *testing.T) {
	assert := assert.New(t)
	resp := httptest.NewRecorder()
	originalData := randomString(1024)
	c := cod.NewContext(resp, nil)

	c.Body = bytes.NewReader([]byte(originalData))

	g := new(GzipCompressor)
	err := g.Pipe(c, 0)
	assert.Nil(err)
	r, err := gzip.NewReader(resp.Body)
	assert.Nil(err)
	defer r.Close()
	buf, _ := ioutil.ReadAll(r)
	assert.Equal(string(buf), originalData)
}
