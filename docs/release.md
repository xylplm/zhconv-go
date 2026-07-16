# 发布文档

## 自动发布（推荐）

仓库已配置 GitHub Actions：

| Actions 显示名 | 文件 | 什么时候跑 | 干什么 | 能否手动 |
|---|---|---|---|---|
| **代码测试** | `ci.yml` | 推代码 / 提 PR / 手动 | 检查 + 单测 + 演示 + 基准冒烟 | ✅ |
| **词典同步** | `dictionary.yml` | 每周一 / 手动 | 从 OpenCC 更新词表，有变化才开 PR（见 [dict-sync.md](dict-sync.md)） | ✅ |
| **版本发布** | `release.yml` | 推送 `v*` tag / 手动填已有 tag | 打多平台 CLI 包并发 GitHub Release | ✅（需已有标签） |

### 手动触发怎么点

1. 打开仓库 **Actions**
2. 左侧点具体工作流（如「代码测试」），不要只停在总览
3. 右侧 **Run workflow** → 选分支 `main` → 绿色运行

说明：

- 按钮只出现在**已启用 `workflow_dispatch` 且文件在默认分支**的工作流页
- **版本发布**手动时必须填**已经 push 过的标签**（如 `v0.1.2`），不会自动帮你新建 tag

### 发布步骤

```bash
# 1. 确保 main 干净且测试通过
go test ./...

# 2. 打 tag（语义化版本）
git tag v0.1.0
git push origin v0.1.0
```

Actions 会自动：

1. 跑测试  
2. 交叉编译：
   - linux/amd64
   - linux/arm64
   - darwin/amd64
   - darwin/arm64
   - windows/amd64
   - windows/arm64  
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
