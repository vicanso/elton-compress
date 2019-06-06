package compress

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/cod"
)

func TestSnappyCompress(t *testing.T) {
	assert := assert.New(t)
	s := new(SnappyCompressor)

	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", snappyEncoding)
	c := cod.NewContext(nil, req)
	acceptable, encoding := s.Accept(c)
	assert.True(acceptable)
	assert.Equal(encoding, snappyEncoding)

	buf, err := s.Compress([]byte("abcd"), 0)
	assert.Nil(err)
	assert.NotEmpty(buf)
}
