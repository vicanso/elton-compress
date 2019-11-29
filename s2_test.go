package compress

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/klauspost/compress/s2"
	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

func TestS2Compress(t *testing.T) {
	assert := assert.New(t)
	originalData := randomString(1024)
	s := new(S2Compressor)

	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", S2Encoding)
	c := elton.NewContext(nil, req)
	acceptable, encoding := s.Accept(c)
	assert.True(acceptable)
	assert.Equal(encoding, S2Encoding)

	buf, err := s.Compress([]byte(originalData), 0)
	assert.Nil(err)
	assert.NotEmpty(buf)
	var originalBuf []byte
	originalBuf, err = s2.Decode(originalBuf, buf)
	assert.Nil(err)
	assert.Equal(originalData, string(originalBuf))
}

func TestS2Pipe(t *testing.T) {
	assert := assert.New(t)
	resp := httptest.NewRecorder()
	originalData := randomString(1024)
	c := elton.NewContext(resp, nil)

	c.Body = bytes.NewReader([]byte(originalData))

	s := new(S2Compressor)
	err := s.Pipe(c, 2)
	assert.Nil(err)
	assert.NotEmpty(resp.Body.Bytes())

	r := s2.NewReader(resp.Body)
	buf, _ := ioutil.ReadAll(r)
	assert.Equal(originalData, string(buf))
}
