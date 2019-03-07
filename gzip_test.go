package compress

import (
	"net/http/httptest"
	"testing"

	"github.com/vicanso/cod"
)

func TestGzipCompress(t *testing.T) {
	g := new(GzipCompressor)
	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	c := cod.NewContext(nil, req)
	acceptable, encoding := g.Accept(c)
	if !acceptable || encoding != gzipEncoding {
		t.Fatalf("request should accept gzip")
	}
	buf, err := g.Compress([]byte("abcd"), 0)
	if err != nil || len(buf) == 0 {
		t.Fatalf("gzip fail, %v", err)
	}
}
