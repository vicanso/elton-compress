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
	Lz4Compressor struct{}
)

// Accept check accept encoding
func (*Lz4Compressor) Accept(c *elton.Context) (acceptable bool, encoding string) {
	return AcceptEncoding(c, Lz4Encoding)
}

// Compress lz4 compress
func (*Lz4Compressor) Compress(buf []byte, level int) ([]byte, error) {
	var b bytes.Buffer
	w := lz4.NewWriter(&b)
	w.Header.CompressionLevel = level
	_, err := w.Write(buf)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Pipe lz4 pipe compress
func (*Lz4Compressor) Pipe(c *elton.Context, level int) (err error) {
	r := c.Body.(io.Reader)
	closer, ok := c.Body.(io.Closer)
	if ok {
		defer closer.Close()
	}
	w := lz4.NewWriter(c.Response)
	w.Header.CompressionLevel = level
	defer w.Close()
	_, err = io.Copy(w, r)
	return
}
