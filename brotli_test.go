// +build brotli

package compress

import (
	"net/http/httptest"
	"testing"

	"github.com/vicanso/cod"
)

func TestBrotliCompress(t *testing.T) {
	br := new(BrCompressor)
	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	c := cod.NewContext(nil, req)
	acceptable, encoding := br.Accept(c)
	if !acceptable || encoding != brEncoding {
		t.Fatalf("request should accept br")
	}
	buf, err := br.Compress([]byte("abcd"), 0)
	if err != nil || len(buf) == 0 {
		t.Fatalf("br fail, %v", err)
	}
}
