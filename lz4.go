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
	"io"

	"github.com/pierrec/lz4"
	"github.com/vicanso/elton"
)

const (
	// Lz4Encoding lz4 encoding
	Lz4Encoding = "lz4"
)

type (
	// Lz4Compressor lz4 compress
	Lz4Compressor struct {
		Level     int
		MinLength int
	}
)

func (l *Lz4Compressor) getMinLength() int {
	if l.MinLength == 0 {
		return defaultCompressMinLength
	}
	return l.MinLength
}

// Accept check accept encoding
func (l *Lz4Compressor) Accept(c *elton.Context, bodySize int) (acceptable bool, encoding string) {
	// 如果数据少于最低压缩长度，则不压缩
	if bodySize >= 0 && bodySize < l.getMinLength() {
		return
	}
	return AcceptEncoding(c, Lz4Encoding)
}

// Compress lz4 compress
func (l *Lz4Compressor) Compress(buf []byte) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	w := lz4.NewWriter(buffer)
	defer w.Close()
	w.Header.CompressionLevel = l.Level
	_, err := w.Write(buf)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

// Pipe lz4 pipe compress
func (l *Lz4Compressor) Pipe(c *elton.Context) (err error) {
	r := c.Body.(io.Reader)
	closer, ok := c.Body.(io.Closer)
	if ok {
		defer closer.Close()
	}
	w := lz4.NewWriter(c.Response)
	w.Header.CompressionLevel = l.Level
	defer w.Close()
	_, err = io.Copy(w, r)
	return
}
