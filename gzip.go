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

	"github.com/vicanso/cod"
)

const (
	gzipEncoding = "gzip"
)

type (
	// GzipCompressor gzip compress
	GzipCompressor struct{}
)

// Accept accept gzip encoding
func (g *GzipCompressor) Accept(c *cod.Context) (acceptable bool, encoding string) {
	return AcceptEncoding(c, gzipEncoding)
}

// Compress compress data by gzip
func (g *GzipCompressor) Compress(buf []byte, level int) ([]byte, error) {
	return doGzip(buf, level)
}

// doGzip 对数据压缩
func doGzip(buf []byte, level int) ([]byte, error) {
	var b bytes.Buffer
	if level <= 0 {
		level = gzip.DefaultCompression
	}
	if level > gzip.BestCompression {
		level = gzip.BestCompression
	}
	w, _ := gzip.NewWriterLevel(&b, level)
	_, err := w.Write(buf)
	if err != nil {
		return nil, err
	}
	w.Close()
	return b.Bytes(), nil
}
