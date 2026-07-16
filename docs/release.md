# 发布文档

## 工作流一览

| Actions 显示名 | 文件 | 什么时候跑 | 干什么 | 能否手动 |
|---|---|---|---|---|
| **代码测试** | `ci.yml` | 推代码 / 提 PR / 手动 | 检查 + 单测 + 演示 + 基准冒烟 | ✅ |
| **词典同步** | `dictionary.yml` | 每周一 / 手动 | 从 OpenCC 更新词表，有变化才开 PR（见 [dict-sync.md](dict-sync.md)） | ✅ |
| **版本发布** | `release.yml` | 推送 `v*` tag / 手动填标签 | 打多平台 CLI 包并发 GitHub Release | ✅ |

## 发布方式

### 方式 A：本地打标签（推荐日常）

```bash
# 1. 确保 main 干净且测试通过
go test ./...

# 2. 打 tag 并推送（语义化版本）
git tag v0.1.2
git push origin v0.1.2
```

推送 `v*` 标签后，**版本发布** 自动跑完测试、交叉编译并创建 Release。

### 方式 B：Actions 手动填写标签

1. 打开 **Actions** → 左侧 **版本发布**
2. **Run workflow**
3. 选分支（一般是 `main`）
4. 填写标签，例如 `v0.1.2`
5. 运行

行为：

| 情况 | 结果 |
|---|---|
| 标签**不存在** | 在所选分支 **当前 HEAD** 创建 annotated tag 并推送，再构建发版 |
| 标签**已存在** | 检出该标签对应提交，再构建发版（可补发/重跑资产） |

说明：

- 标签格式：`v主.次.修订`，可选后缀如 `v0.2.0-rc.1`
- 手动新建标签时由 `github-actions[bot]` 推送；为避免再触发一轮 tag push 发版，bot 推送的 tag 事件会被跳过
- 不会改写已存在的标签指向（已有 tag 只用于构建，不 force-move）

### 手动触发其它工作流

1. 打开仓库 **Actions**
2. 左侧点具体工作流，不要只停在总览
3. **Run workflow** → 选分支 `main` → 运行

按钮只出现在**已启用 `workflow_dispatch` 且文件在默认分支**的工作流页。

## 发布流水线做什么

1. 跑 `go vet` + `go test`  
2. 交叉编译 CLI：
   - linux/amd64、linux/arm64  
   - darwin/amd64、darwin/arm64  
   - windows/amd64、windows/arm64  
3. 打包 `.tar.gz` / `.zip`  
4. 生成 `checksums.txt`  
5. 创建 GitHub Release 并上传资产  

### 资产命名

```text
zhconv_<version>_<os>_<arch>.tar.gz
zhconv_<version>_windows_<arch>.zip
checksums.txt
```

## 本地打包

```bash
./scripts/build-release.sh v0.1.0
ls dist/
```

Windows 可用 Git Bash 执行。

## go install

```bash
go install github.com/xylplm/zhconv-go/cmd/zhconv@v0.1.0
zhconv -version
zhconv -demo
```

## 库版本引用

```bash
go get github.com/xylplm/zhconv-go@v0.1.0
```
