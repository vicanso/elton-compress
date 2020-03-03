package compress

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/golang/snappy"
	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

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

	r := snappy.NewReader(resp.Body)
	buf, _ := ioutil.ReadAll(r)
	assert.Equal(originalData, string(buf))
}
