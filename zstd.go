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

	"github.com/klauspost/compress/zstd"
	"github.com/vicanso/elton"
)

const (
	// ZstdEncoding zstd encoding
	ZstdEncoding = "zst"
)

type (
	// ZstdCompressor zstd compress
	ZstdCompressor struct{}
)

func getZstdEncoderLevel(level int) zstd.EncoderLevel {
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
func (*ZstdCompressor) Compress(buf []byte, level int) ([]byte, error) {
	encoder, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(getZstdEncoderLevel(level)))
	if err != nil {
		return nil, err
	}
	return encoder.EncodeAll(buf, make([]byte, 0, len(buf))), nil
}

// Pipe zstd pike
func (*ZstdCompressor) Pipe(c *elton.Context, level int) (err error) {
	r := c.Body.(io.Reader)
	closer, ok := c.Body.(io.Closer)
	if ok {
		defer closer.Close()
	}
	w, err := zstd.NewWriter(c.Response, zstd.WithEncoderLevel(getZstdEncoderLevel(level)))
	if err != nil {
		return err
	}

	defer w.Close()
	_, err = io.Copy(w, r)
	return
}
