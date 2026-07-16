# zhconv-go

纯 Go 实现的**繁体 → 简体**中文转换库。

- 单向：只做 `t2s`（面向 zh-Hans）
- 词组优先最长匹配，单字兜底
- 覆盖：通用繁体字形 + 台湾用语词组 + 港台异体字反向映射
- 内嵌词表，零 CGO
- 并发安全、边界安全（非法 UTF-8 透传）
- 可 `go get github.com/xylplm/zhconv-go`

词表来源于 [OpenCC](https://github.com/BYVoid/OpenCC)（Apache-2.0）：
`TSCharacters` / `TSPhrases` / `TWPhrasesRev`，以及 `TWVariants` / `HKVariants` 的反向字形映射。

> 说明：这是**单向繁转简**，不是完整 OpenCC/zhconv 多地区互转引擎。  
> 地区支持体现在“台/港繁体写法也能落到大陆简体”，而不是 `zh-TW/zh-HK` 目标变体切换。

## 安装

```bash
go get github.com/xylplm/zhconv-go@latest
```

## 使用

```go
package main

import (
	"fmt"

	"github.com/xylplm/zhconv-go"
)

func main() {
	fmt.Println(zhconv.ToSimplified("軟體與網路連線"))
	// 软件与网络连线
}
```

自定义转换器：

```go
c, err := zhconv.New(zhconv.Options{})
if err != nil {
	panic(err)
}
out := c.Convert("資料庫")
```

仅单字、禁用词组：

```go
c, _ := zhconv.New(zhconv.Options{DisablePhrases: true})
```

## 命令行试玩

```bash
go run ./cmd/zhconv -demo

echo 軟體與網路 | go run ./cmd/zhconv
go run ./cmd/zhconv -i in.txt -o out.txt
```

## 设计

```
phrases trie (longest match)
        ↓ miss
character map (rune -> simplified)
        ↓ miss
original text
```

- 词表构建一次，查询热路径避免多余分配
- `Default()` 使用 `sync.Once` 共享实例
- API 面积极小：`ToSimplified` / `Converter.Convert`

## 词表

| 文件 | 说明 |
|---|---|
| `table/chars.tsv` | 繁→简 单字（含港台异体反向） |
| `table/phrases.tsv` | 繁→简 词组（含台湾用语） |
| `dict/NOTICE` | 上游许可说明 |

## 地区覆盖现状

| 类型 | 支持程度 |
|---|---|
| 通用繁体字形 | ✅ 字符表全量 t2s |
| 台湾用语（软體/網路/伺服器…） | ✅ 词组表 |
| 港台异体字（裏/裡、啓/啟、綫/線…） | ✅ 字符变体反向 |
| `zh-TW`/`zh-HK` 作为输出目标 | ❌ 不做（只输出简体） |
| 简体 → 繁体 | ❌ 不做 |

## 测试

```bash
go test ./...
go test -bench=BenchmarkToSimplified -benchmem ./...
```

## License

Apache-2.0（词表遵循 OpenCC 的 Apache-2.0）
