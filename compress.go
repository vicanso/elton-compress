// Copyright 2018 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compress

import (
	"regexp"
	"strings"

	"github.com/vicanso/elton"
)

var (
	defaultCompressRegexp = regexp.MustCompile("text|javascript|json")
)

const (
	defaultCompressMinLength = 1024
)

type (
	// Compressor compressor interface
	Compressor interface {
		// Accept accept check function
		Accept(c *elton.Context) (acceptable bool, encoding string)
		// Compress compress function
		Compress([]byte, int) ([]byte, error)
		// Pipe pipe function
		Pipe(*elton.Context, int) error
	}
	// Config compress config
	Config struct {
		// Level compress level
		Level int
		// MinLength min compress length
		MinLength int
		// Checker check the data is compressable
		Checker *regexp.Regexp
		// CompressorList compressor list
		CompressorList []Compressor
		// Skipper skipper function
		Skipper elton.Skipper
	}
)

// AcceptEncoding check request accept encoding
func AcceptEncoding(c *elton.Context, encoding string) (bool, string) {
	acceptEncoding := c.GetRequestHeader(elton.HeaderAcceptEncoding)
	if strings.Contains(acceptEncoding, encoding) {
		return true, encoding
	}
	return false, ""
}

// NewDefault create a default compress middleware, support gzip
func NewDefault() elton.Handler {
	return NewWithDefaultCompressor(Config{})
}

// NewWithDefaultCompressor create compress middleware with default compressor
func NewWithDefaultCompressor(config Config) elton.Handler {
	compressorList := make([]Compressor, 0)

	// 添加默认的 brotli 压缩
	br := new(BrCompressor)
	_, err := br.Compress([]byte("brotli"), 0)
	// 如果可以压缩成功，则添加 br 压缩
	if err == nil {
		compressorList = append(compressorList, br)
	}

	// 添加默认的 gzip 压缩
	compressorList = append(compressorList, new(GzipCompressor))
	config.CompressorList = compressorList

	return New(config)
}

// New create a new compress middleware
func New(config Config) elton.Handler {
	minLength := config.MinLength
	if minLength == 0 {
		minLength = defaultCompressMinLength
	}
	skipper := config.Skipper
	if skipper == nil {
		skipper = elton.DefaultSkipper
	}
	checker := config.Checker
	if checker == nil {
		checker = defaultCompressRegexp
	}
	compressorList := config.CompressorList
	return func(c *elton.Context) (err error) {
		if skipper(c) || compressorList == nil {
			return c.Next()
		}
		err = c.Next()
		if err != nil {
			return
		}
		isReaderBody := c.IsReaderBody()
		// 如果数据为空，而且body不是reader，直接跳过
		if c.BodyBuffer == nil && !isReaderBody {
			return
		}

		// encoding 不为空，已做处理，无需要压缩
		if c.GetHeader(elton.HeaderContentEncoding) != "" {
			return
		}
		contentType := c.GetHeader(elton.HeaderContentType)
		// 数据类型为非可压缩，则返回
		if !checker.MatchString(contentType) {
			return
		}

		var body []byte
		if c.BodyBuffer != nil {
			body = c.BodyBuffer.Bytes()
		}
		if !isReaderBody {
			// 如果数据长度少于最小压缩长度
			if len(body) < minLength {
				return
			}
		}

		for _, compressor := range compressorList {
			acceptable, encoding := compressor.Accept(c)
			if !acceptable {
				continue
			}
			if isReaderBody {
				c.SetHeader(elton.HeaderContentEncoding, encoding)
				err = compressor.Pipe(c, config.Level)
				if err != nil {
					return
				}
				// pipe 将数据直接转至原有的Response，因此设置committed为true
				c.Committed = true
				// 清除 reader body
				c.Body = nil
			} else {
				newBuf, e := compressor.Compress(body, config.Level)
				// 如果压缩成功，则使用压缩数据
				// 失败则忽略
				if e == nil {
					c.SetHeader(elton.HeaderContentEncoding, encoding)
					c.BodyBuffer.Reset()
					c.BodyBuffer.Write(newBuf)
					break
				}
			}
		}
		return
	}
}
