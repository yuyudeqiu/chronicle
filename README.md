# Chronicle

Chronicle 这是一个专为 LLM Agent 打造的**任务管理与追踪系统 (Agent-Centric Task Management System)**。
系统提供了一组稳定、原子化的后台 API，接管用户的日常任务拆解、进度追踪与日志记录，通过 CLI 和 HTTP API 进行管理。

核心设计原则：
1. **防止幻觉：** 通过极简的列表接口让 Agent 获取真实 ID；通过复合接口保证状态更新与日志记录的原子性（数据库事务）。
2. **自动化归档：** 提供机器友好的 JSON 日报接口（供 Agent 写总结），以及 SQLite 数据库导入导出命令，方便备份和迁移。

## ✨ 核心特性

- 🤖 **Agent-Friendly API**: 专门设计的防止大模型产生幻觉的接口格式，严格限制 Token 消耗。
- 🔄 **原子化事务**: 状态流转和日志记录确保数据库层面的强一致性，避免脏数据。
- ⌨️ **CLI**: 在终端管理任务、工作日志、归档、日报和周报，支持 JSON 输出，无需启动服务。
- 📦 **轻量级**: 使用 Go 和 SQLite 构建，仅需单一可执行文件即可运行，无需复杂的环境依赖。
- 📤 **数据导入导出**: 支持 `export` / `import` 命令，方便数据备份与迁移。

## 🛠 技术栈

- **后端**: Golang, Gin 框架, GORM
- **数据库**: SQLite (本地单文件存储, 纯 Go 驱动)

## 📂 项目结构

```text
.
├── cmd/                          # 命令行工具逻辑
│   ├── root.go                   # 根命令定义
│   ├── server.go                 # 服务器启动命令
│   ├── tasks.go                  # 任务管理相关命令
│   ├── archive.go                # 归档、恢复归档和删除工作日志
│   ├── version.go                # 版本信息命令
│   └── export.go                 # 数据导入导出命令
├── internal/                     # 核心业务逻辑
│   ├── config/                   # 配置管理（支持环境变量和命令行参数）
│   ├── handler/                  # HTTP 请求处理与响应 (Controller)
│   ├── model/                    # GORM 实体与接口 DTO
│   └── service/                  # 核心业务逻辑 (事务处理、状态机等)
├── skills/                       # AI Agent 集成 Skill 文件目录
├── .github/workflows/            # GitHub Actions CI 配置
├── Makefile                      # 构建 Makefile
├── DESIGN.md                     # 技术设计文档
├── main.go                       # 项目主入口
└── go.mod                        # Go 模块定义
```

## 🚀 快速开始

### 1. 环境准备

请确保您本地已经安装了：
- [Golang](https://go.dev/dl/) (版本要求 >= 1.24)

### 2. 安装 Chronicle

#### 方式一：go install（推荐）

```bash
go install github.com/yuyudeqiu/chronicle@latest
```

安装 `main` 分支的开发版（版本信息会显示 commit）：

```bash
go install github.com/yuyudeqiu/chronicle@main
```

#### 方式二：手动编译

```bash
# 使用 Makefile（推荐）
make build

# 或直接使用 go build
go build -o chronicle main.go
```

### 3. 数据存储配置

默认情况下，数据存放在当前工作目录下的 `data/app.db`，首次运行会自动创建数据目录。

可通过以下方式自定义数据存储路径（优先级：命令行参数 > 环境变量 > 默认值）：

```bash
# 方式一：环境变量
export CHRONICLE_DATA_DIR=/custom/path

# 方式二：命令行参数
chronicle --data-dir /custom/path create "新任务"
```

### 4. 启动 HTTP API（可选）

```bash
# 启动 HTTP API 服务器，默认端口为 8080
chronicle server

# 或使用环境变量指定数据目录
CHRONICLE_DATA_DIR=/path/to/data chronicle server
```

服务只提供 JSON API，不再托管网页。可以通过版本接口检查服务：

```bash
curl http://localhost:8080/api/v1/version
```

API 路由见下表。CLI 直接访问 SQLite，与 API 共用业务逻辑；使用 CLI 无需先启动服务。
如需读写同一份数据，请为两者配置相同的 `--data-dir` 或 `CHRONICLE_DATA_DIR`。

### 5. 使用命令行工具 (CLI)

直接在终端管理任务：

```bash
# 查看所有命令
chronicle --help

# 创建任务
chronicle create "完成项目重构" -c 开发 -d "使用 Cobra 优化 CLI"

# 列出活跃任务（待办、进行中和阻塞）
chronicle list

# 列出已完成任务（最近 50 条）
chronicle list done

# 查看任务详情
chronicle get <task_id>

# 添加执行日志
chronicle log <task_id> "完成了 CLI 重构"

# 更新任务状态
chronicle update <task_id> --new-status done

# 更新任务信息（包括标题）
chronicle update <task_id> --title "新的任务标题" -c 运维 -d "更新描述" --deadline 2026-01-01T00:00:00Z

# 归档、查看归档任务、恢复归档
chronicle archive <task_id>
chronicle list archived
chronicle unarchive <task_id>

# 删除指定工作日志（get 会显示日志 ID，也可使用 get --json）
chronicle delete-worklog <worklog_id>

# 删除任务
chronicle delete <task_id>

# 查看每日总结
chronicle summary
chronicle summary 2026-01-15

# 查看周报
chronicle weekly-summary

# 查看任务统计
chronicle stats

# 查看版本信息
chronicle version
chronicle version --json

# 导出数据库
chronicle export /path/to/backup.db

# 导入数据库（会覆盖现有数据）
chronicle import /path/to/backup.db
chronicle import /path/to/backup.db --force  # 跳过确认
```

#### 全局参数

| 参数 | 说明 |
|------|------|
| `--data-dir <path>` | 指定数据目录（默认 `data/`，也可通过 `CHRONICLE_DATA_DIR` 环境变量设置） |
| `-o, --json` | 以 JSON 格式输出结果（便于 Agent 解析） |

### HTTP API

所有接口使用 `/api/v1` 前缀，返回 JSON；未知路径返回 JSON 404。

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | `/version` | 版本信息 |
| POST | `/tasks` | 创建任务 |
| GET | `/tasks` | 活跃任务；`?status=done` 获取已完成任务 |
| GET | `/tasks/archived` | 归档任务 |
| GET / PATCH / DELETE | `/tasks/:id` | 查询、修改、删除任务 |
| POST | `/tasks/:id/progress` | 添加工作日志并更新进度 |
| POST | `/tasks/:id/archive` | 归档任务 |
| POST | `/tasks/:id/unarchive` | 恢复归档 |
| DELETE | `/worklogs/:id` | 删除工作日志 |
| GET | `/reports/daily-summary` | 日报，可传 `?date=YYYY-MM-DD` |
| GET | `/reports/weekly-summary` | 周报，可传 `?date=YYYY-MM-DD` |
| GET | `/stats/summary` | 任务统计 |

### 📚 AI Agent 集成

如果你是 AI Agent 想接入 Chronicle，推荐使用 `skills/` 目录下的 Skill 示例代码。

**快速集成示例：**

```bash
# 列出任务（JSON 格式，便于解析）
chronicle list --json

# 创建任务
chronicle create "任务标题" -c "分类"

# 添加记录
chronicle log <id> "工作内容"

# 标记完成
chronicle update <id> --new-status done
```

详细文档见：[skills/](./skills/) 目录

## 发布与版本规范

Chronicle 使用 `v0.x.x` 语义化版本。合入 `main` 后，Release Please 会持续更新
Release PR；只有合并该 Release PR 才会创建版本 tag 和 GitHub Release。完整的版本
规则、提交格式及首次启用步骤见 [RELEASING.md](RELEASING.md)。
