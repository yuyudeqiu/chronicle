# Chronicle 技术设计

## 架构

Chronicle 保留 CLI 和 HTTP API 两个入口，共用 `internal/service` 中的业务逻辑和 SQLite 数据库。
CLI 直接调用 service，无需启动 HTTP 服务。`chronicle server` 仅提供 `/api/v1` JSON API。
项目不包含 Web 前端，不依赖 Node.js 或静态资源，也未实现 TUI。

```text
cmd/                     Cobra 命令：任务、归档、日志、报表、导入导出、服务启动
internal/handler/        Gin HTTP 路由和请求响应处理
internal/service/        任务操作、事务、归档和报表
internal/model/          数据模型、请求和响应结构
internal/config/         数据目录配置
internal/buildinfo/      CLI 和 API 共用的版本信息
```

## 数据与配置

- SQLite 默认路径为 `data/app.db`，首次运行会创建数据目录并迁移表结构。
- 数据目录优先级：`--data-dir` > `CHRONICLE_DATA_DIR` > `data/`。
- 默认路径相对于进程工作目录；CLI 和 API 使用同一份数据时应指定相同的绝对路径。
- `tasks` 保存任务信息，状态支持 `todo`、`in-progress`、`blocked`、`done`。
- `task_logs` 通过 `task_id` 关联任务，保存工作日志。
- `archived_at` 标记归档，归档不删除任务或日志；活跃列表和已完成列表排除归档任务。
- 状态变为 `done` 时记录完成时间；进度接口在一个事务中添加日志并更新任务。
- 删除任务时，在一个事务中先删除日志再删除任务。

## 接口与命令

HTTP 接口及 CLI 示例见 [使用参考](docs/usage.md)，实际路由定义位于
`internal/handler/task.go`。成功响应使用 `{"code":0,"msg":"success","data":...}`，
未知路径返回 HTTP 404 和 JSON 错误。

列表返回任务 ID、标题、分类、状态等精简字段，详情接口返回完整任务和日志。
CLI 支持 `--json`，便于 Agent 和脚本解析；`get` 的文本输出同时显示工作日志 ID。

CLI 覆盖任务创建、查询、更新、删除、归档、恢复归档、添加和删除工作日志、统计、日报和周报。
`export` / `import` 用于 SQLite 数据库文件备份和迁移。
目前未实现 Markdown 导出接口。

## 构建与部署

- `make build` 只编译当前平台的 Go 可执行文件，包含 CLI 和 API server。
- GoReleaser 构建 macOS、Linux 和 Windows 发布包，无需额外静态资源。
