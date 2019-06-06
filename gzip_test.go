package compress

import (
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
	g := new(GzipCompressor)
	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	c := cod.NewContext(nil, req)
	acceptable, encoding := g.Accept(c)
	assert.True(acceptable)
	assert.Equal(encoding, gzipEncoding)
	buf, err := g.Compress([]byte("abcd"), 0)
	assert.Nil(err)
	assert.NotEmpty(buf)
}
