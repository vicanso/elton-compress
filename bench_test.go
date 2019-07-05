package compress

import (
	"fmt"
	"io/ioutil"
	"testing"
)

var (
	benchData []byte
)

func init() {
	buf, err := ioutil.ReadFile("./assets/index.html")
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
