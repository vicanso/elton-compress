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
	"math/rand"
	"testing"
)

// benchData 为本地生成的测试语料。基准测试不依赖网络，
// 否则离线环境下整个包的 go test 都会失败
var benchData = newBenchData(64 * 1024)

// newBenchData 生成指定大小的伪文本语料：由固定词表按确定的随机序列拼接而成，
// 兼具真实文本的可压缩性与跨机器可复现性
func newBenchData(size int) []byte {
	words := bytes.Fields(compressTestData)
	// 固定种子，保证每次运行的语料完全一致
	r := rand.New(rand.NewSource(1))
	buf := new(bytes.Buffer)
	buf.Grow(size)
	for buf.Len() < size {
		buf.Write(words[r.Intn(len(words))])
		buf.WriteByte(' ')
	}
	return buf.Bytes()[:size]
}

func BenchmarkSnappy(b *testing.B) {
	benchmarkCompress(b, NewSnappyCompressor().Compress)
}

func BenchmarkSnappyBest(b *testing.B) {
	benchmarkCompress(b, (&SnappyCompressor{Level: SnappyLevelBest}).Compress)
}

func BenchmarkS2(b *testing.B) {
	benchmarkCompress(b, NewS2Compressor().Compress)
}

func BenchmarkS2Best(b *testing.B) {
	benchmarkCompress(b, (&S2Compressor{Level: S2LevelBest}).Compress)
}

func BenchmarkLz4(b *testing.B) {
	benchmarkCompress(b, NewLz4Compressor().Compress)
}

func BenchmarkLz4Best(b *testing.B) {
	benchmarkCompress(b, (&Lz4Compressor{Level: Lz4LevelBest}).Compress)
}

func benchmarkCompress(b *testing.B, fn func([]byte, ...int) (*bytes.Buffer, error)) {
	b.SetBytes(int64(len(benchData)))
	b.ReportAllocs()
	b.ResetTimer()
	var buf *bytes.Buffer
	for i := 0; i < b.N; i++ {
		var err error
		buf, err = fn(benchData)
		if err != nil {
			b.Fatal(err)
		}
	}
	// 压缩率随基准一起输出，便于横向对比各算法。
	// ResetTimer 会清空已上报的指标，因此必须在循环之后上报
	b.ReportMetric(float64(len(benchData))/float64(buf.Len()), "ratio")
}
