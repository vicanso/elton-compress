// +build brotli

package compress

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/cod"
)

func TestBrotliCompress(t *testing.T) {
	assert := assert.New(t)
	br := new(BrCompressor)
	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	c := cod.NewContext(nil, req)
	acceptable, encoding := br.Accept(c)
	assert.True(acceptable)
	assert.Equal(encoding, brEncoding)
	buf, err := br.Compress([]byte("abcd"), 0)
	assert.Nil(err)
	assert.NotEqual(len(buf), 0)
}
