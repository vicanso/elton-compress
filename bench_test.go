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
	fmt.Println(fmt.Sprintf("original: %d bytes", len(buf)))
	g := new(GzipCompressor)
	gzipBuf, _ := g.Compress(buf, 0)
	fmt.Println(fmt.Sprintf("gzip: %d bytes", len(gzipBuf)))

	br := new(BrCompressor)
	brBuf, _ := br.Compress(buf, 0)
	fmt.Println(fmt.Sprintf("br: %d bytes", len(brBuf)))

	sn := new(SnappyCompressor)
	snappyBuf, _ := sn.Compress(buf, 0)
	fmt.Println(fmt.Sprintf("snappy: %d bytes", len(snappyBuf)))

	s2 := new(S2Compressor)
	s2Buf, _ := s2.Compress(buf, 0)
	fmt.Println(fmt.Sprintf("s2: %d bytes", len(s2Buf)))

	z := new(ZstdCompressor)
	zBuf, _ := z.Compress(buf, 0)
	fmt.Println(fmt.Sprintf("zstd: %d bytes", len(zBuf)))

	lz := new(Lz4Compressor)
	lzBuf, _ := lz.Compress(buf, 0)
	fmt.Println(fmt.Sprintf("lz4: %d bytes", len(lzBuf)))
}

func BenchmarkGzip(b *testing.B) {
	b.ReportAllocs()
	g := new(GzipCompressor)
	for i := 0; i < b.N; i++ {
		_, err := g.Compress(benchData, 0)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkBr(b *testing.B) {
	b.ReportAllocs()
	br := new(BrCompressor)
	for i := 0; i < b.N; i++ {
		_, err := br.Compress(benchData, 0)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkSnappy(b *testing.B) {
	b.ReportAllocs()
	sn := new(SnappyCompressor)
	for i := 0; i < b.N; i++ {
		_, err := sn.Compress(benchData, 0)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkS2(b *testing.B) {
	b.ReportAllocs()
	s2 := new(S2Compressor)
	for i := 0; i < b.N; i++ {
		_, err := s2.Compress(benchData, 0)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkS2Fast(b *testing.B) {
	b.ReportAllocs()
	s2 := new(S2Compressor)
	for i := 0; i < b.N; i++ {
		_, err := s2.Compress(benchData, 1)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkZstd(b *testing.B) {
	b.ReportAllocs()
	z := new(ZstdCompressor)
	for i := 0; i < b.N; i++ {
		_, err := z.Compress(benchData, 0)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkLz4(b *testing.B) {
	b.ReportAllocs()
	lz := new(Lz4Compressor)
	for i := 0; i < b.N; i++ {
		_, err := lz.Compress(benchData, 0)
		if err != nil {
			panic(err)
		}
	}
}
