# Chronicle

面向 AI Agent 的轻量任务管理工具。通过 CLI 或 HTTP API 管理任务、记录工作日志、生成日报和周报。
使用 Go 和 SQLite，数据保存在本地，CLI 无需启动服务。

## 安装

需要 Go 1.24 或更高版本：

```bash
go install github.com/yuyudeqiu/chronicle@latest
```

## 快速开始

```bash
# 创建任务并查看 ID
chronicle create "完成项目重构" -c 开发
chronicle list

# 将 <task_id> 替换为上面返回的 ID
chronicle log <task_id> "完成了 CLI 重构"
chronicle update <task_id> --new-status done

# 查看今日工作记录
chronicle summary
```

添加 `--json` 可输出 JSON，便于 Agent 和脚本解析：

```bash
chronicle list --json
```

默认数据库位于**当前工作目录**下的 `data/app.db`，首次运行自动创建。
跨目录使用时，建议设置固定的绝对路径：

```bash
export CHRONICLE_DATA_DIR="$HOME/.chronicle"
```

也可通过 `--data-dir /path/to/data` 覆盖。

## 更多用法

- 查看命令：`chronicle --help` 或 `chronicle <command> --help`
- 启动 HTTP API：`chronicle server`（默认端口 8080）
- [使用参考](docs/usage.md)：完整命令示例、配置、备份、API 路由与源码构建
- [Agent Skill](skills/chronicle/SKILL.md)：任务管理、每日计划和工作总结
- [技术设计](DESIGN.md) · [发布规范](RELEASING.md) · [更新日志](CHANGELOG.md)
