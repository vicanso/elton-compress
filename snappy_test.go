package compress

import (
	"net/http/httptest"
	"testing"

	"github.com/vicanso/cod"
)

func TestSnappyCompress(t *testing.T) {
	s := new(SnappyCompressor)

	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", snappyEncoding)
	c := cod.NewContext(nil, req)
	acceptable, encoding := s.Accept(c)
	if !acceptable || encoding != snappyEncoding {
		t.Fatalf("request should accept snappy")
	}
	buf, err := s.Compress([]byte("abcd"), 0)
	if err != nil || len(buf) == 0 {
		t.Fatalf("snappy fail, %v", err)
	}
}
