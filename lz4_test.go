package compress

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/pierrec/lz4"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

func TestLz4Compress(t *testing.T) {
	assert := assert.New(t)
	originalData := randomString(1024)
	z := new(Lz4Compressor)

	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", Lz4Encoding)
	c := elton.NewContext(nil, req)
	acceptable, encoding := z.Accept(c, 0)
	assert.False(acceptable)
	assert.Empty(encoding)
	acceptable, encoding = z.Accept(c, len(originalData))
	assert.True(acceptable)
	assert.Equal(encoding, Lz4Encoding)

	buf, err := z.Compress([]byte(originalData))
	assert.Nil(err)
	assert.NotEmpty(buf)

	r := lz4.NewReader(buf)
	dst, err := ioutil.ReadAll(r)
	assert.Nil(err)
	assert.Equal([]byte(originalData), dst)
}

func TestLz4Pipe(t *testing.T) {
	assert := assert.New(t)
	resp := httptest.NewRecorder()
	originalData := randomString(1024)
	c := elton.NewContext(resp, nil)
	c.Body = bytes.NewReader([]byte(originalData))

	z := new(Lz4Compressor)
	err := z.Pipe(c)
	assert.Nil(err)
	assert.NotEmpty(resp.Body.Bytes())

	r := lz4.NewReader(resp.Body)
	dst, err := ioutil.ReadAll(r)
	assert.Nil(err)
	assert.Equal([]byte(originalData), dst)
}
