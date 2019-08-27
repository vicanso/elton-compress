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
// limitations under the License

package compress

import (
	"bytes"
	"io"

	"github.com/andybalholm/brotli"
	"github.com/vicanso/elton"
)

const (
	brEncoding       = "br"
	maxBrQuality     = 11
	defaultBrQuality = 6
)

type (
	// BrCompressor brotli compress
	BrCompressor struct{}
)

func getBrLevel(level int) int {
	if level <= 0 {
		level = defaultBrQuality
	}
	if level > maxBrQuality {
		level = maxBrQuality
	}
	return level
}

// Accept check accept econding
func (b *BrCompressor) Accept(c *elton.Context) (acceptable bool, encoding string) {
	return AcceptEncoding(c, brEncoding)
}

// Compress brotli compress
func (b *BrCompressor) Compress(buf []byte, level int) ([]byte, error) {
	buffer := new(bytes.Buffer)
	w := brotli.NewWriterLevel(buffer, getBrLevel(level))
	_, err := w.Write(buf)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Pipe brotli pipe
func (b *BrCompressor) Pipe(c *elton.Context, level int) (err error) {
	r := c.Body.(io.Reader)
	closer, ok := c.Body.(io.Closer)
	if ok {
		defer closer.Close()
	}
	w := brotli.NewWriterLevel(c.Response, getBrLevel(level))

	defer w.Close()
	_, err = io.Copy(w, r)
	return
}
