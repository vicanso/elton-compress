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

// +build brotli

package compress

import (
	"io"

	"github.com/google/brotli/go/cbrotli"
	"github.com/vicanso/cod"
)

const (
	brEncoding = "br"
	maxQuality = 11
)

type (
	// BrCompressor brotli compress
	BrCompressor struct{}
)

func getBrLevel(level int) int {
	if level <= 0 {
		level = 9
	}
	if level > maxQuality {
		level = maxQuality
	}
	return level
}

// Accept check accept econding
func (b *BrCompressor) Accept(c *cod.Context) (acceptable bool, encoding string) {
	return AcceptEncoding(c, brEncoding)
}

// Compress brotli compress
func (b *BrCompressor) Compress(buf []byte, level int) ([]byte, error) {
	return cbrotli.Encode(buf, cbrotli.WriterOptions{
		Quality: getBrLevel(level),
		LGWin:   0,
	})
}

// Pipe brotli pipe
func (b *BrCompressor) Pipe(c *cod.Context, level int) (err error) {
	r := c.Body.(io.Reader)
	closer, ok := c.Body.(io.Closer)
	if ok {
		defer closer.Close()
	}
	w := cbrotli.NewWriter(c.Response, cbrotli.WriterOptions{
		Quality: getBrLevel(level),
		LGWin:   0,
	})
	defer w.Close()
	_, err = io.Copy(w, r)
	return
}
