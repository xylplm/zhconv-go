# 使用文档

## 1. 作为 Go 库

### 安装

```bash
go get github.com/xylplm/zhconv-go@latest
```

`go.mod`：

```go
require github.com/xylplm/zhconv-go v0.1.0
```

### 基本转换

```go
package main

import (
	"fmt"

	"github.com/xylplm/zhconv-go"
)

func main() {
	in := "請將繁體中文字幕轉成簡體"
	out := zhconv.ToSimplified(in)
	fmt.Println(out)
	// 请将繁体中文字幕转成简体
}
```

### 高频场景：复用 Converter

`ToSimplified` 内部使用共享实例，已经足够。  
若你需要自定义选项，或希望显式持有实例：

```go
c, err := zhconv.New(zhconv.Options{})
if err != nil {
	panic(err)
}

// Converter 构造完成后可并发调用 Convert。
out1 := c.Convert("軟體")
out2 := c.Convert("網路")
_ = out1
_ = out2
```

### 仅单字映射（禁用词组）

```go
c, err := zhconv.New(zhconv.Options{DisablePhrases: true})
if err != nil {
	panic(err)
}
// “軟體”会按单字变成“软体”，不一定是词组结果“软件”
fmt.Println(c.Convert("軟體"))
```

### 处理文件

```go
package main

import (
	"os"

	"github.com/xylplm/zhconv-go"
)

func main() {
	data, err := os.ReadFile("traditional.txt")
	if err != nil {
		panic(err)
	}
	out := zhconv.ToSimplifiedBytes(data)
	if err := os.WriteFile("simplified.txt", out, 0o644); err != nil {
		panic(err)
	}
}
```

### 并发

```go
c := zhconv.Default()
// 多 goroutine 同时 c.Convert(...) 是安全的
```

### 边界行为

| 输入 | 行为 |
|---|---|
| `""` | 返回 `""` |
| 已是简体 | 基本保持不变 |
| ASCII / 数字 / 标点 | 保持不变 |
| 非法 UTF-8 字节 | 原样透传，不 panic |
| `(*Converter)(nil).Convert(s)` | 返回原字符串 |

## 2. 命令行 CLI

### 安装

```bash
go install github.com/xylplm/zhconv-go/cmd/zhconv@latest
```

或下载 Release 资产，例如：

- `zhconv_Linux_x86_64.tar.gz`
- `zhconv_Windows_x86_64.zip`
- `zhconv_Darwin_arm64.tar.gz`

### 演示

```bash
zhconv -demo
```

输出示例：

```text
IN : 軟體與網路連線
OUT: 软件与网络连接
```

### 标准输入

```bash
echo 資料庫程式設計師 | zhconv
# 数据库程序员
```

### 文件互转

```bash
zhconv -i traditional.srt.txt -o simplified.srt.txt
```

> 当前 CLI/库都是**纯文本层**转换。  
> 若用于 ASS/SRT 字幕文件，请先自行保证只送入对白文本，或在上层做格式安全切分（后续可扩展）。

## 3. 常见转换样例

| 输入 | 输出 |
|---|---|
| 軟體與網路連線 | 软件与网络连接 |
| 資料庫伺服器 | 数据库服务器 |
| 螢幕上的訊息 | 屏幕上的消息 |
| 檔案已匯出 | 文件已导出 |
| 裏面/裡面 | 里面 |
| 啓動/啟動 | 启动 |

## 4. 错误与性能建议

1. **不要每个请求 `New` 一次**  
   用 `Default()` 或进程内单例。

2. **大文本**  
   转换是线性扫描；字幕/文档级别通常毫秒级。

3. **质量**  
   这是规则转换，不是翻译模型。专有名词偶发不准时，靠扩充词组表解决。

4. **与业务集成建议**  
   - 先判定是否繁体主导  
   - 再 `ToSimplified`  
   - 再做业务侧校验（例如字幕语言复验）

## 5. 版本与兼容

- 模块路径：`github.com/xylplm/zhconv-go`
- 语义化版本：`vMajor.Minor.Patch`
- `v0.x` 允许在必要时刻调整 API，但会尽量保持 `ToSimplified` / `Convert` 稳定
