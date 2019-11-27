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
	"io"

	"github.com/golang/snappy"
	"github.com/vicanso/elton"
)

const (
	snappyEncoding = "snz"
)

type (
	// SnappyCompressor snappy compress
	SnappyCompressor struct{}
)

// Accept check accept encoding
func (s *SnappyCompressor) Accept(c *elton.Context) (acceptable bool, encoding string) {
	return AcceptEncoding(c, snappyEncoding)
}

// Compress snappy compress
func (s *SnappyCompressor) Compress(buf []byte, level int) ([]byte, error) {
	var dst []byte
	data := snappy.Encode(dst, buf)
	return data, nil
}

// Pipe snappy pipe
func (s *SnappyCompressor) Pipe(c *elton.Context, level int) (err error) {
	r := c.Body.(io.Reader)
	closer, ok := c.Body.(io.Closer)
	if ok {
		defer closer.Close()
	}
	w := snappy.NewBufferedWriter(c.Response)
	defer w.Close()
	_, err = io.Copy(w, r)
	return
}
