# elton-compress 2.0 迁移指南

2.0 跟随 [elton 2.0](https://github.com/vicanso/elton/blob/master/docs/migration-v2.md) 升级，模块路径变更为 `github.com/vicanso/elton-compress/v2`，最低要求 **Go 1.24**。

```bash
go get github.com/vicanso/elton-compress/v2
```

```go
import (
	compress "github.com/vicanso/elton-compress/v2"
	"github.com/vicanso/elton/v2"
	"github.com/vicanso/elton/v2/middleware"
)
```

`Compressor` 接口（`Accept` / `Compress` / `Pipe`）本身未变，压缩器仍然直接传给 `middleware.NewCompressConfig`。

## 移除 ZstdCompressor

elton 2.0 的 middleware 包已内置 zstd，本模块不再重复提供。

```go
// v1
config := middleware.NewCompressConfig(new(compress.ZstdCompressor))

// v2
config := middleware.NewCompressConfig(middleware.NewZstdCompressor())
```

同时移除 `compress.ZstdEncoding`，改用 `elton.Zstd`。缓存场景可用 `middleware.NewCacheZstdCompressor()`。

内置版本顺带修复了本模块 v1 中的一个缺陷：`ZstdCompressor.Compress` 每次调用都会通过 `zstd.NewWriter(nil, ...)` 创建 encoder 却从不 `Close()`，而 encoder 会按 GOMAXPROCS 启动 goroutine，导致每次压缩都泄漏一批 goroutine。

## LZ4 改为标准 frame 格式

v1 输出的是裸 LZ4 block（`lz4.CompressBlock`），**不携带原始数据长度**，客户端必须通过带外方式得知解压后的大小才能解码，作为 HTTP `Content-Encoding` 无法互操作。

v2 改为输出标准 LZ4 frame（`lz4.NewWriter`），任意 lz4 客户端可直接解压，`Pipe` 也变成真正的流式压缩而不再全量读入内存。

这是字节流层面的破坏性变更：**升级服务端后，按 v1 裸 block 实现的客户端需要同步改为 frame 解码器。**

依赖同时从已废弃的 `github.com/pierrec/lz4`（v2.6.1+incompatible）升级到 `github.com/pierrec/lz4/v4`。

## Lz4Compressor.Level 开始生效

v1 的 `Level` 字段与 `Compress` 的 level 参数被完全忽略（`lz4.CompressBlock` 硬编码为最快档），README 中"级别越高压缩率越好"的描述并不成立。

v2 中 `Level` 取值为 `1`（`Lz4LevelFastest`，对应 `lz4.Fast`）到 `10`（`Lz4LevelBest`，对应 `lz4.Level9`），默认为 1。**若此前依赖"设了 Level 也不会变慢"的行为，升级后 CPU 开销会随 Level 上升。**

## 新增 S2Compressor

v1 的 README 中已出现 `compress.S2Compressor` 示例，但代码中并不存在该类型。v2 补齐了这一实现：编码为 `s2`，级别 1-3（默认 1）。

s2 是 snappy 的超集，同级别下速度与压缩率均优于 snappy，但产生的数据只能由 s2 解码器解压。

## Snappy 级别语义修正

`SnappyEncoding` 仍为 `"snz"`，block 格式不变，已部署客户端无需改动。

v1 中默认路径走 `snappy.Encode`（在当前 klauspost/compress 中等价于 `s2.EncodeSnappyBetter`），而显式传 `level=1` 走的是更弱的 `s2.EncodeSnappy`——即"指定级别 1"反而比"不指定"压得更差。

v2 用 `SnappyLevelFastest`(1) / `SnappyLevelBetter`(2) / `SnappyLevelBest`(3) 明确三档，零值 `Level` 的默认档位为 `SnappyLevelBetter`，与 v1 的默认输出保持一致。

## 其它

- 新增构造函数 `NewSnappyCompressor()` / `NewS2Compressor()` / `NewLz4Compressor()`，与 elton 2.0 的压缩器风格一致。
- 新增包级压缩/解压函数：`SnappyCompress`/`SnappyDecompress`、`S2Compress`/`S2Decompress`、`Lz4Compress`/`Lz4Decompress`，可脱离中间件单独使用。
- 三个压缩器统一处理 `middleware.IgnoreCompression` 与越界级别：非正值回退到各自默认档，超出上限收敛至上限，不再返回错误。
- 移除 `ErrLz4IsNotCompressible`。frame 格式对不可压缩数据会自动退化为存储块，不再需要这一错误。
- `Pipe` 在压缩出错时不再吞掉 `Close()` 的错误（frame 结束标记的写入失败会被返回）。
- 基准测试不再于 `init()` 中下载 jquery，改用本地生成的语料，离线环境下 `go test` 不会再直接 panic。
