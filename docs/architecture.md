# 架构说明

## 目标

做一个**轻量、高效、可维护**的单向繁转简库：

- 只输出简体（zh-Hans 取向）
- 输入可覆盖通用繁体、台湾用语、港台异体
- 纯 Go，零 CGO
- API 极简，适合被其他项目 `go get`

## 包结构

```text
zhconv-go/
  converter.go       # Converter + trie 最长匹配
  default.go         # Default/ToSimplified 单例
  doc.go             # 包文档
  table/
    embed.go         # go:embed 词表
    table.go         # TSV 解析
    chars.tsv
    phrases.tsv
  cmd/zhconv/        # CLI
  examples/basic/    # 最小可运行示例
  docs/              # 文档
  .github/workflows/ # CI / Release
```

## 转换流水线

```text
input string
  │
  ├─ phrase trie longest match  ──命中──► write simplified phrase
  │
  └─ miss
       │
       ├─ rune char map ──命中──► write simplified char(s)
       │
       └─ miss ──► write original
```

### 为什么词组优先

很多繁简差异不能只靠单字，例如：

- 軟體 → 软件（不是“软体”）
- 記憶體 → 内存
- 伺服器 → 服务器

### 为什么词组目标还要再过单字表

OpenCC 风格数据常是“链式”：

1. 地区词先落到某一中间繁/开式写法  
2. 再靠字符表落到最终简体  

因此加载词组时会对 `to` 再执行一次字符级归一，避免结果残留繁体。

## 数据结构

### 字符表

```go
char1 map[rune]rune   // 热路径 1:1
charN map[rune]string // 极少数一对多
```

- 绝大多数映射落在 `char1`，避免每字产生临时 string
- `charN` 仅承载目标为多码点的稀有项

### 词组表

```go
phrases map[rune][]phrase // 按首字分桶，桶内按词长降序
```

- 从左到右扫描，首字入桶后按最长优先比较
- 匹配过程零临时 `[]rune` 分配，直接对输入做 UTF-8 比较
- 词表规模小（约 1k 词组），比通用 AC 自动机更易维护

### 零分配快路径

- 输入无需转换时：`Convert` 直接返回原字符串（0 alloc）
- `ConvertBytes` 无变化时返回原 `[]byte`

## 并发模型

- `Converter` 在构建完成后只读
- `Default()` 用 `sync.Once` 初始化一次
- 多 goroutine 同时 `Convert` 安全

## 词表来源与归一

| 源 | 作用 |
|---|---|
| `TSCharacters` | 通用繁→简单字 |
| `TSPhrases` | 通用繁→简词组 |
| `TWPhrasesRev` | 台湾用语 → 简体侧 |
| `TWVariants` / `HKVariants` | 反向：地区异体 → 开放字形 → 简体 |

归一规则：

- 多候选取第一个
- `from == to` 丢弃
- 非法/空行忽略

## 扩展点

后续若要增强，优先这些低耦合扩展：

1. **字幕格式适配层**（SRT/ASS 只转对白）  
2. **更大词组包**（可选 `phrases_ext.tsv`）  
3. **自定义词表**（`Options.Chars/Phrases` 已支持）  
4. **质量统计**（命中词组数/字符数）

不建议扩展成：

- 多目标变体引擎（zh-TW/HK/CN 全矩阵）
- 运行时远程拉词表
- CGO / 外部二进制依赖

## 非目标

- 简转繁
- 翻译润色
- 语义消歧（超出规则转换）
