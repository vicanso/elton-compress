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

// Package compress 为 elton 的 compress 中间件补充 elton 未内置的压缩算法：
// snappy、s2 与 lz4，适用于对 CPU 敏感、压缩率要求不高的内网场景。
//
// gzip、br、zstd 已由 github.com/vicanso/elton/v2/middleware 内置，
// 请直接使用 middleware.NewGzipCompressor / NewBrCompressor / NewZstdCompressor。
package compress

import (
	"io"

	"github.com/vicanso/elton/v2"
)

// normalizeLevel 把压缩级别收敛到 [1, maxLevel] 区间。
// 中间件在未配置 DynamicLevel 时不传级别，压缩器零值 Level 也是 0，
// 同时 middleware.IgnoreCompression 为负数，这些非正值统一回退到 defaultLevel。
func normalizeLevel(level, defaultLevel, maxLevel int) int {
	if level <= 0 {
		return defaultLevel
	}
	return min(level, maxLevel)
}

// pipeReader 取出 c.Body 中的 reader，并在返回的函数中关闭它（如果它实现了 io.Closer）。
// 中间件仅在 c.IsReaderBody() 为真时才调用 Pipe，因此此处的断言是安全的。
func pipeReader(c *elton.Context) (io.Reader, func()) {
	r := c.Body.(io.Reader)
	closer, ok := c.Body.(io.Closer)
	if !ok {
		return r, func() {}
	}
	return r, func() {
		_ = closer.Close()
	}
}
