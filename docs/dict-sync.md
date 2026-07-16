# 词典自动同步

词表来自 [OpenCC](https://github.com/BYVoid/OpenCC)（Apache-2.0），本仓库只维护**单向繁→简**产物。

## 生成规则

| 上游文件 | 用途 |
|---|---|
| `TSCharacters.txt` | 主字符表（第一候选） |
| `TSPhrases.txt` | 基础词组 |
| `TWPhrasesRev.txt` | 台湾用语反向词组（同 key 覆盖基础词组） |
| `TWVariants.txt` / `HKVariants.txt` | 地区异体反向补洞（不覆盖已有 TS 字符） |

生成入口：

```bash
go run ./scripts/gendict
# 或
./scripts/sync-dict.sh
# 指定上游 ref
OPENCC_REF=master go run ./scripts/gendict
go run ./scripts/gendict -ref v1.1.9
```

输出：

```text
dict/chars.tsv
dict/phrases.tsv
dict/SOURCE.json      # 上游 commit / 计数 / 时间（仅词表内容变化时更新）
dict/NOTICE
table/chars.tsv       # go:embed 副本，必须与 dict 一致
table/phrases.tsv
```

无变更行为：

- TSV 映射内容与仓库内一致时，`gendict` **跳过写入**（含 `SOURCE.json`）
- 避免仅因 `generated_at` 每周开空 PR

## GitHub Actions 自动 PR

工作流：`.github/workflows/dictionary.yml`  
Actions 列表显示名：**词典同步**

| 触发 | 行为 |
|---|---|
| 每周一 03:17 UTC（定时） | 拉取 OpenCC 最新词表，有变化才开/更新 PR |
| 手动 Run workflow | 可选手动指定 OpenCC 版本 |

PR 特性：

- 固定分支：`chore/dictionary-sync`（有更新则刷新同一 PR）
- 标题：`chore 🔧: 同步 OpenCC 词典 <short-sha>`
- 合并前已跑检查与测试
- 只提交词表相关文件

### 手动触发

GitHub → Actions → **词典同步** → Run workflow

### 权限说明

默认 `GITHUB_TOKEN` 即可在本仓库开 PR。  
若仓库启用了更严的 branch protection / 需要跨 fork，再考虑换成 PAT。

## 本地验收建议

```bash
./scripts/sync-dict.sh
go test ./... -count=1
git diff --stat -- dict table
```

关注：

1. 字符/词组数量是否异常暴跌  
2. 已知用例是否仍通过（`軟體`、`伺服器`、`裡面`…）  
3. `dict` 与 `table` 是否完全一致  

## 与发布的关系

词典 PR **不会**自动打 release tag。  
合并后如需让下游拿到新词表，再按 [release.md](release.md) 发版。
