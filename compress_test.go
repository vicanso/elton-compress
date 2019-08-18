package compress

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")

type testCompressor struct{}

func (t *testCompressor) Accept(c *elton.Context) (acceptable bool, encoding string) {
	return AcceptEncoding(c, "br")
}

func (t *testCompressor) Compress(buf []byte, level int) ([]byte, error) {
	return []byte("abcd"), nil
}

func (t *testCompressor) Pipe(c *elton.Context, level int) error {
	return nil
}

// randomString get random string
func randomString(n int) string {
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestAcceptEncoding(t *testing.T) {
	assert := assert.New(t)
	req := httptest.NewRequest("GET", "/", nil)
	c := elton.NewContext(nil, req)
	acceptable, encoding := AcceptEncoding(c, elton.Gzip)
	assert.False(acceptable)
	assert.Empty(encoding)

	c.SetRequestHeader(elton.HeaderAcceptEncoding, elton.Gzip)
	acceptable, encoding = AcceptEncoding(c, elton.Gzip)
	assert.True(acceptable)
	assert.Equal(elton.Gzip, encoding)
}

func TestCompress(t *testing.T) {
	t.Run("skip", func(t *testing.T) {
		assert := assert.New(t)
		c := elton.NewContext(nil, nil)
		done := false
		c.Next = func() error {
			done = true
			return nil
		}
		fn := New(Config{
			Skipper: func(c *elton.Context) bool {
				return true
			},
		})
		err := fn(c)
		assert.Nil(err)
		assert.True(done)
	})

	t.Run("nil body", func(t *testing.T) {
		assert := assert.New(t)
		c := elton.NewContext(nil, nil)
		done := false
		c.Next = func() error {
			done = true
			return nil
		}
		fn := NewDefault()
		err := fn(c)
		assert.Nil(err)
		assert.True(done)
	})

	t.Run("return error", func(t *testing.T) {
		assert := assert.New(t)
		c := elton.NewContext(nil, nil)
		customErr := errors.New("abccd")
		c.Next = func() error {
			return customErr
		}
		fn := NewDefault()
		err := fn(c)
		assert.Equal(err, customErr)
	})

	t.Run("normal", func(t *testing.T) {
		assert := assert.New(t)
		conf := Config{
			Level:     1,
			MinLength: 1,
		}
		conf.AddCompressor(new(GzipCompressor))
		fn := New(conf)

		req := httptest.NewRequest("GET", "/users/me", nil)
		req.Header.Set(elton.HeaderAcceptEncoding, "gzip")
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		c.SetHeader(elton.HeaderContentType, "text/html")
		c.BodyBuffer = bytes.NewBuffer([]byte("<html><body>" + randomString(8192) + "</body></html>"))
		originalSize := c.BodyBuffer.Len()
		done := false
		c.Next = func() error {
			done = true
			return nil
		}
		err := fn(c)
		assert.Nil(err)
		assert.True(done)
		assert.True(c.BodyBuffer.Len() < originalSize)
		assert.Equal(elton.Gzip, c.GetHeader(elton.HeaderContentEncoding))
	})

	t.Run("encoding done", func(t *testing.T) {
		assert := assert.New(t)
		fn := NewDefault()
		req := httptest.NewRequest("GET", "/users/me", nil)
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		c.Next = func() error {
			return nil
		}
		body := bytes.NewBufferString(randomString(4096))
		c.BodyBuffer = body
		c.SetHeader(elton.HeaderContentEncoding, "gzip")
		err := fn(c)
		assert.Nil(err)
		assert.Equal(body.Bytes(), c.BodyBuffer.Bytes())
	})

	t.Run("body size is less than min length", func(t *testing.T) {
		assert := assert.New(t)
		fn := NewDefault()

		req := httptest.NewRequest("GET", "/users/me", nil)
		req.Header.Set(elton.HeaderAcceptEncoding, "gzip")
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		c.Next = func() error {
			return nil
		}
		body := bytes.NewBufferString("abcd")
		c.BodyBuffer = body
		c.SetHeader(elton.HeaderContentType, "text/plain")
		err := fn(c)
		assert.Nil(err)
		assert.Equal(body.Bytes(), c.BodyBuffer.Bytes())
		assert.Empty(c.GetHeader(elton.HeaderContentEncoding))
	})

	t.Run("image should not be compress", func(t *testing.T) {
		assert := assert.New(t)

		fn := NewDefault()

		req := httptest.NewRequest("GET", "/users/me", nil)
		req.Header.Set(elton.HeaderAcceptEncoding, "gzip")
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		c.SetHeader(elton.HeaderContentType, "image/jpeg")
		c.Next = func() error {
			return nil
		}
		body := bytes.NewBufferString(randomString(4096))
		c.BodyBuffer = body
		err := fn(c)
		assert.Nil(err)
		assert.Equal(body.Bytes(), c.BodyBuffer.Bytes())
		assert.Empty(c.GetHeader(elton.HeaderContentEncoding))
	})

	t.Run("not accept gzip should not compress", func(t *testing.T) {
		assert := assert.New(t)

		fn := NewDefault()

		req := httptest.NewRequest("GET", "/users/me", nil)
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		c.SetHeader(elton.HeaderContentType, "text/html")
		c.Next = func() error {
			return nil
		}
		body := bytes.NewBufferString(randomString(4096))
		c.BodyBuffer = body
		err := fn(c)
		assert.Nil(err)
		assert.Equal(body.Bytes(), c.BodyBuffer.Bytes())
		assert.Empty(c.GetHeader(elton.HeaderContentEncoding))
	})

	t.Run("custom compress", func(t *testing.T) {
		assert := assert.New(t)
		compressorList := make([]Compressor, 0)
		compressorList = append(compressorList, new(testCompressor))
		fn := New(Config{
			Compressors: compressorList,
		})

		req := httptest.NewRequest("GET", "/users/me", nil)
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		c.SetHeader(elton.HeaderContentType, "text/html")
		c.BodyBuffer = bytes.NewBufferString("<html><body>" + randomString(8192) + "</body></html>")
		done := false
		c.Next = func() error {
			done = true
			return nil
		}
		err := fn(c)
		assert.Nil(err)
		assert.True(done)
		assert.Equal(4, c.BodyBuffer.Len())
		assert.Equal("br", c.GetHeader(elton.HeaderContentEncoding))
	})

	t.Run("reader body", func(t *testing.T) {
		assert := assert.New(t)

		fn := NewDefault()

		req := httptest.NewRequest("GET", "/users/me", nil)
		req.Header.Set(elton.HeaderAcceptEncoding, "gzip")
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		c.SetHeader(elton.HeaderContentType, "text/html")
		c.Next = func() error {
			return nil
		}
		body := bytes.NewBufferString(randomString(4096))
		gzipBytes, _ := doGzip(body.Bytes(), 0)
		c.Body = body
		err := fn(c)
		assert.True(c.Committed)
		assert.Nil(err)
		assert.Equal(gzipBytes, resp.Body.Bytes())
		assert.Equal(elton.Gzip, c.GetHeader(elton.HeaderContentEncoding))
	})
}

// https://stackoverflow.com/questions/50120427/fail-unit-tests-if-coverage-is-below-certain-percentage
func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.85 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
	os.Exit(rc)
}
