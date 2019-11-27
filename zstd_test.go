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
	req.Header.Set("Accept-Encoding", zstdEncoding)
	c := elton.NewContext(nil, req)
	acceptable, encoding := z.Accept(c)
	assert.True(acceptable)
	assert.Equal(encoding, zstdEncoding)

	buf, err := z.Compress([]byte(originalData), 0)
	assert.Nil(err)
	assert.NotEmpty(buf)

	decorder, err := zstd.NewReader(nil)
	assert.Nil(err)
	var dst []byte
	dst, err = decorder.DecodeAll(buf, dst)
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
	err := z.Pipe(c, 9)
	assert.Nil(err)
	assert.NotEmpty(resp.Body.Bytes())

	decorder, err := zstd.NewReader(nil)
	assert.Nil(err)
	var dst []byte
	dst, err = decorder.DecodeAll(resp.Body.Bytes(), dst)
	assert.Nil(err)
	assert.Equal([]byte(originalData), dst)
}
