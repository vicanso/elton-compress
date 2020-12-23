// MIT License

// Copyright (c) 2020 Tree Xie

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package compress

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/pierrec/lz4"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

var compressTestData = []byte(`Brotli is a data format specification[2] for data streams compressed with a specific combination of the general-purpose LZ77 lossless compression algorithm, Huffman coding and 2nd order context modelling. Brotli is a compression algorithm developed by Google and works best for text compression.

Google employees Jyrki Alakuijala and Zoltán Szabadka initially developed Brotli to decrease the size of transmissions of WOFF2 web fonts, and in that context Brotli was a continuation of the development of zopfli, which is a zlib-compatible implementation of the standard gzip and deflate specifications. Brotli allows a denser packing than gzip and deflate because of several algorithmic and format level improvements: the use of context models for literals and copy distances, describing copy distances through past distances, use of move-to-front queue in entropy code selection, joint-entropy coding of literal and copy lengths, the use of graph algorithms in block splitting, and a larger backward reference window are example improvements. The Brotli specification was generalized in September 2015 for HTTP stream compression (content-encoding type 'br'). This generalized iteration also improved the compression ratio by using a pre-defined dictionary of frequently used words and phrases.`)

func doLZ4Decode(buf []byte) ([]byte, error) {
	dst := make([]byte, 10*len(buf))
	n, err := lz4.UncompressBlock(buf, dst)
	if err != nil {
		return nil, err
	}
	dst = dst[:n]
	return dst, nil
}

func TestLz4Compress(t *testing.T) {
	assert := assert.New(t)
	originalData := compressTestData
	z := new(Lz4Compressor)

	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Accept-Encoding", Lz4Encoding)
	c := elton.NewContext(nil, req)
	acceptable, encoding := z.Accept(c, 0)
	assert.False(acceptable)
	assert.Empty(encoding)
	acceptable, encoding = z.Accept(c, len(originalData))
	assert.True(acceptable)
	assert.Equal(encoding, Lz4Encoding)

	buf, err := z.Compress([]byte(originalData))
	assert.Nil(err)
	assert.NotEmpty(buf)

	dst, err := doLZ4Decode(buf.Bytes())
	assert.Nil(err)
	assert.Equal([]byte(originalData), dst)
}

func TestLz4Pipe(t *testing.T) {
	assert := assert.New(t)
	resp := httptest.NewRecorder()
	originalData := compressTestData
	c := elton.NewContext(resp, nil)
	c.Body = bytes.NewReader([]byte(originalData))

	z := new(Lz4Compressor)
	err := z.Pipe(c)
	assert.Nil(err)
	assert.NotEmpty(resp.Body.Bytes())

	dst, err := doLZ4Decode(resp.Body.Bytes())
	assert.Nil(err)
	assert.Equal([]byte(originalData), dst)
}
