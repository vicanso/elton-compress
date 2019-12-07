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

	"github.com/klauspost/compress/zstd"
	"github.com/vicanso/elton"
)

const (
	// ZstdEncoding zstd encoding
	ZstdEncoding = "zst"
)

type (
	// ZstdCompressor zstd compress
	ZstdCompressor struct {
		Level int
	}
)

func (z *ZstdCompressor) getLevel() zstd.EncoderLevel {
	level := z.Level
	l := zstd.EncoderLevel(level)
	if l < zstd.SpeedFastest || l > zstd.SpeedBestCompression {
		return zstd.SpeedDefault
	}
	return l
}

// Accept check accept encoding
func (*ZstdCompressor) Accept(c *elton.Context) (acceptable bool, encoding string) {
	return AcceptEncoding(c, ZstdEncoding)
}

// Compress zstd compress
func (z *ZstdCompressor) Compress(buf []byte) (*bytes.Buffer, error) {
	encoder, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(z.getLevel()))
	if err != nil {
		return nil, err
	}
	data := encoder.EncodeAll(buf, make([]byte, 0, len(buf)))
	return bytes.NewBuffer(data), nil
}

// Pipe zstd pike
func (z *ZstdCompressor) Pipe(c *elton.Context) (err error) {
	r := c.Body.(io.Reader)
	closer, ok := c.Body.(io.Closer)
	if ok {
		defer closer.Close()
	}
	w, err := zstd.NewWriter(c.Response, zstd.WithEncoderLevel(z.getLevel()))
	if err != nil {
		return err
	}

	defer w.Close()
	_, err = io.Copy(w, r)
	return
}
