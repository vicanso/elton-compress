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
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

var (
	benchData []byte
)

func init() {
	resp, err := http.Get("https://code.jquery.com/jquery-3.4.1.min.js")
	if err != nil {
		panic(err)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	benchData = buf

	br := new(BrCompressor)
	brBuf, _ := br.Compress(buf)
	fmt.Println(fmt.Sprintf("br: %d bytes", brBuf.Len()))

	sn := new(SnappyCompressor)
	snappyBuf, _ := sn.Compress(buf)
	fmt.Println(fmt.Sprintf("snappy: %d bytes", snappyBuf.Len()))

	s2 := new(S2Compressor)
	s2Buf, _ := s2.Compress(buf)
	fmt.Println(fmt.Sprintf("s2: %d bytes", s2Buf.Len()))

	z := new(ZstdCompressor)
	zBuf, _ := z.Compress(buf)
	fmt.Println(fmt.Sprintf("zstd: %d bytes", zBuf.Len()))

	lz := new(Lz4Compressor)
	lzBuf, _ := lz.Compress(buf)
	fmt.Println(fmt.Sprintf("lz4: %d bytes", lzBuf.Len()))
}

func BenchmarkBr(b *testing.B) {
	b.ReportAllocs()
	br := new(BrCompressor)
	for i := 0; i < b.N; i++ {
		_, err := br.Compress(benchData)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkSnappy(b *testing.B) {
	b.ReportAllocs()
	sn := new(SnappyCompressor)
	for i := 0; i < b.N; i++ {
		_, err := sn.Compress(benchData)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkS2(b *testing.B) {
	b.ReportAllocs()
	s2 := new(S2Compressor)
	for i := 0; i < b.N; i++ {
		_, err := s2.Compress(benchData)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkS2Fast(b *testing.B) {
	b.ReportAllocs()
	s2 := new(S2Compressor)
	s2.Level = 1
	for i := 0; i < b.N; i++ {
		_, err := s2.Compress(benchData)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkZstd(b *testing.B) {
	b.ReportAllocs()
	z := new(ZstdCompressor)
	for i := 0; i < b.N; i++ {
		_, err := z.Compress(benchData)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkLz4(b *testing.B) {
	b.ReportAllocs()
	lz := new(Lz4Compressor)
	for i := 0; i < b.N; i++ {
		_, err := lz.Compress(benchData)
		if err != nil {
			panic(err)
		}
	}
}
