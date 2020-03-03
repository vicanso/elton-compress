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
	"bytes"
	"compress/gzip"
	"io"

	"github.com/vicanso/elton"
)

const (
	// GzipEncoding gzip encoding
	GzipEncoding = "gzip"
)

type (
	// GzipCompressor gzip compress
	GzipCompressor struct {
		Level     int
		MinLength int
	}
)

// Accept accept gzip encoding
func (g *GzipCompressor) Accept(c *elton.Context, bodySize int) (acceptable bool, encoding string) {
	// 如果数据少于最低压缩长度，则不压缩
	if bodySize >= 0 && bodySize < g.getMinLength() {
		return
	}
	return AcceptEncoding(c, GzipEncoding)
}

// Compress compress data by gzip
func (g *GzipCompressor) Compress(buf []byte) (*bytes.Buffer, error) {
	level := g.getLevel()
	buffer := new(bytes.Buffer)

	w, _ := gzip.NewWriterLevel(buffer, level)
	defer w.Close()
	_, err := w.Write(buf)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (g *GzipCompressor) getLevel() int {
	level := g.Level
	if level <= 0 {
		level = gzip.DefaultCompression
	}
	if level > gzip.BestCompression {
		level = gzip.BestCompression
	}
	return level
}

func (g *GzipCompressor) getMinLength() int {
	if g.MinLength == 0 {
		return defaultCompressMinLength
	}
	return g.MinLength
}

// Pipe compress by pipe
func (g *GzipCompressor) Pipe(c *elton.Context) (err error) {
	r := c.Body.(io.Reader)
	closer, ok := c.Body.(io.Closer)
	if ok {
		defer closer.Close()
	}
	w, _ := gzip.NewWriterLevel(c.Response, g.getLevel())
	defer w.Close()
	_, err = io.Copy(w, r)
	return
}
