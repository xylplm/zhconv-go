# zhconv-go

纯 Go 实现的**繁体 → 简体**中文转换库。

[![CI](https://github.com/xylplm/zhconv-go/actions/workflows/ci.yml/badge.svg)](https://github.com/xylplm/zhconv-go/actions/workflows/ci.yml)
[![Release](https://github.com/xylplm/zhconv-go/actions/workflows/release.yml/badge.svg)](https://github.com/xylplm/zhconv-go/actions/workflows/release.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/xylplm/zhconv-go.svg)](https://pkg.go.dev/github.com/xylplm/zhconv-go)

- 单向：只做 `t2s`（面向 zh-Hans）
- 词组优先最长匹配，单字兜底
- 覆盖：通用繁体字形 + 台湾用语词组 + 港台异体字反向映射
- 内嵌词表，零 CGO
- 并发安全、边界安全（非法 UTF-8 透传）
- 可 `go get github.com/xylplm/zhconv-go`

词表来源于 [OpenCC](https://github.com/BYVoid/OpenCC)（Apache-2.0）：
`TSCharacters` / `TSPhrases` / `TWPhrasesRev`，以及 `TWVariants` / `HKVariants` 的反向字形映射。

词典可本地重生，也可由 Actions **每周自动开 PR** 同步上游：

```bash
make dict-sync
# 或
go run ./scripts/gendict
```

详见 [docs/dict-sync.md](docs/dict-sync.md)。

> 说明：这是**单向繁转简**，不是完整 OpenCC/zhconv 多地区互转引擎。  
> 地区支持体现在“台/港繁体写法也能落到大陆简体”，而不是 `zh-TW/zh-HK` 目标变体切换。

## 安装

### 作为库

```bash
go get github.com/xylplm/zhconv-go@latest
```

要求：**Go 1.26+**

### 安装 CLI

```bash
go install github.com/xylplm/zhconv-go/cmd/zhconv@latest
```

或从 [Releases](https://github.com/xylplm/zhconv-go/releases) 下载对应平台二进制。

## 30 秒上手

```go
package main

import (
	"fmt"

	"github.com/xylplm/zhconv-go"
)

func main() {
	fmt.Println(zhconv.ToSimplified("軟體與網路連線"))
	// 软件与网络连接
}
```

更多文档与示例：

- [docs/usage.md](docs/usage.md) — 完整用法
- [docs/architecture.md](docs/architecture.md) — 架构设计
- [docs/comparison.md](docs/comparison.md) — 与 OpenCC / zhconv-rs 等主流方案对比
- [docs/dict-sync.md](docs/dict-sync.md) — 词典同步与自动 PR
- [docs/release.md](docs/release.md) — 发布流程
- [examples/basic](examples/basic) — 可运行示例

## 库 API

### 最简

```go
out := zhconv.ToSimplified("資料庫程式設計師")
// 数据库程序员
```

### 共享实例（推荐高频调用）

```go
c := zhconv.Default() // sync.Once，进程内单例
out := c.Convert("伺服器檔案已匯出")
```

### 自定义选项

```go
c, err := zhconv.New(zhconv.Options{
	// DisablePhrases: true, // 仅单字映射
})
if err != nil {
	panic(err)
}
fmt.Println(c.Convert("螢幕解析度"))
```

### 字节接口

```go
out := zhconv.ToSimplifiedBytes([]byte("繁體中文"))
```

## 命令行

```bash
# 内置演示
zhconv -demo

# stdin -> stdout
echo 軟體與網路連線 | zhconv

# 文件
zhconv -i in.txt -o out.txt
```

本地开发：

```bash
go run ./cmd/zhconv -demo
```

## 地区覆盖

| 类型 | 支持程度 |
|---|---|
| 通用繁体字形 | ✅ 字符表全量 t2s |
| 台湾用语（软體/網路/伺服器…） | ✅ 词组表 |
| 港台异体字（裏/裡、啓/啟、綫/線…） | ✅ 字符变体反向 |
| `zh-TW`/`zh-HK` 作为输出目标 | ❌ 不做（只输出简体） |
| 简体 → 繁体 | ❌ 不做 |

## 与主流方案怎么选

| 你的目标 | 更合适 |
|---|---|
| Go 项目里只要繁→简，轻量可嵌入 | **zhconv-go** |
| 完整双向/多地区配置（s2t、t2tw、s2hk…） | OpenCC / zhconv-rs |
| Go 里要较完整 OpenCC 能力 | `longbridgeapp/opencc` 等 |

多维度对比图与详细表格见：**[docs/comparison.md](docs/comparison.md)**。

## 设计摘要

```
phrases trie (longest match)
        ↓ miss
character map (rune -> simplified)
        ↓ miss
original text
```

- 词表构建一次，热路径低分配
- `Default()` 使用 `sync.Once`
- API 面积极小：`ToSimplified` / `Converter.Convert`

详见 [docs/architecture.md](docs/architecture.md)。

## 词表

| 文件 | 说明 |
|---|---|
| `table/chars.tsv` | 繁→简 单字（含港台异体反向） |
| `table/phrases.tsv` | 繁→简 词组（含台湾用语） |
| `dict/NOTICE` | 上游许可说明 |

## 测试

```bash
go test ./...
go test -bench=BenchmarkToSimplified -benchmem .
```

## 本地构建多平台二进制

```bash
# Linux / macOS / Git Bash
./scripts/build-release.sh v0.1.0

# 产物在 dist/
```

## 发布

推送语义化版本 tag 即可触发 GitHub Actions 构建并创建 Release：

```bash
git tag v0.1.0
git push origin v0.1.0
```

工作流：

- CI：`.github/workflows/ci.yml`（测试 + vet）
- Release：`.github/workflows/release.yml`（多平台打包 + checksums + GitHub Release）

## License

Apache-2.0  
词表遵循 OpenCC 的 Apache-2.0（见 `dict/NOTICE`）。
