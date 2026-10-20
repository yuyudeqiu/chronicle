# Chronicle 使用参考

[返回 README](../README.md)

## 数据存储

默认情况下，数据存放在当前工作目录下的 `data/app.db`，首次运行会自动创建数据目录。

可通过以下方式自定义数据存储路径（优先级：命令行参数 > 环境变量 > 默认值）：

```bash
# 方式一：环境变量
export CHRONICLE_DATA_DIR=/custom/path

# 方式二：命令行参数
chronicle --data-dir /custom/path create "新任务"
```

## CLI

示例中的 `<task_id>` 和 `<worklog_id>` 需要替换为实际 ID。
任务状态支持 `todo`、`in-progress`、`blocked`、`done`；截止时间使用 RFC 3339 格式，例如 `2026-10-20T18:00:00+08:00`。

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

# 按单个或多个状态筛选（不包含归档任务）
chronicle list todo
chronicle list todo,in-progress

# 查看任务详情
chronicle get <task_id>

# 添加执行日志
chronicle log <task_id> "完成了 CLI 重构"

# 更新任务状态
chronicle update <task_id> --new-status done

# 同时更新状态和其他信息
chronicle update <task_id> --new-status done --title "最终标题" -d "已完成验证"

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

更新状态与其他信息可在同一次 `update` 中完成。首次进入 `done` 时记录完成时间；重复设置 `done` 保留原完成时间；从 `done` 改回其他状态时清除完成时间。
`update --json` 返回更新后的完整任务，包括状态与完成时间。

### 全局参数

| 参数 | 说明 |
|------|------|
| `--data-dir <path>` | 指定数据目录（默认 `data/`，也可通过 `CHRONICLE_DATA_DIR` 环境变量设置） |
| `-o, --json` | 以 JSON 格式输出结果（便于 Agent 解析） |

## HTTP API

```bash
chronicle server
curl http://localhost:8080/api/v1/version
```

默认端口为 8080。CLI 无需启动服务；两者读写同一份数据时，请配置相同的绝对数据目录。

所有接口使用 `/api/v1` 前缀，返回 JSON；未知路径返回 JSON 404。

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | `/version` | 版本信息 |
| POST | `/tasks` | 创建任务 |
| GET | `/tasks` | 默认活跃任务；`?status=todo` 或 `?status=todo,in-progress` 按状态筛选；`?status=done` 获取最近 50 条已完成任务 |
| GET | `/tasks/archived` | 归档任务 |
| GET / PATCH / DELETE | `/tasks/:id` | 查询、修改、删除任务 |
| POST | `/tasks/:id/progress` | 添加工作日志并更新进度 |
| POST | `/tasks/:id/archive` | 归档任务 |
| POST | `/tasks/:id/unarchive` | 恢复归档 |
| DELETE | `/worklogs/:id` | 删除工作日志 |
| GET | `/reports/daily-summary` | 日报，可传 `?date=YYYY-MM-DD` |
| GET | `/reports/weekly-summary` | 周报，可传 `?date=YYYY-MM-DD` |
| GET | `/stats/summary` | 任务统计 |

CLI 与 API 的状态筛选支持 `todo`、`in-progress`、`blocked`、`done`，多个状态以逗号分隔，顺序无关，重复值自动去重。筛选排除归档任务；单独筛选 `done` 返回最近 50 条，其他组合按截止时间和创建时间排序，不限制条数。非法状态或空的组合项会报错，API 返回 HTTP 400。

添加日志并将任务标记为完成时，向 `POST /tasks/:id/progress` 传入：

```json
{"log_text":"完成实现并验证","new_status":"done"}
```

完成状态统一使用 `new_status: "done"`，不再支持 `mark_as_done` 参数。

## 从源码构建

在仓库根目录运行：

```bash
make build
./bin/chronicle --help
```

也可使用 `go build -o chronicle main.go`。Windows 下可执行文件带 `.exe` 后缀。
安装 main 分支的开发版可使用 `go install github.com/yuyudeqiu/chronicle@main`，版本信息会显示 commit。
