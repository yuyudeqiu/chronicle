---
name: chronicle
description: 使用 Chronicle CLI 管理任务、记录工作进度，并根据本地任务数据生成每日计划或工作总结。
---

# Chronicle

使用已安装的 `chronicle` 命令；CLI 直接读写本地 SQLite，无需启动 HTTP 服务。
沿用用户配置的 `CHRONICLE_DATA_DIR` 或 `--data-dir`；默认 `data/` 相对于当前工作目录，跨目录运行时应使用同一个绝对数据路径。

## 操作约定

- 查询时使用 `--json`。更新已有任务前，通过 `list --json` 或 `get <task_id> --json` 获取真实 ID 和当前信息；有多个候选时先确认目标。
- `list` 返回待办、进行中和阻塞任务；已完成任务用 `list done`，归档任务用 `list archived`。
- 状态支持 `todo`、`in-progress`、`blocked`、`done`。添加日志不会自动改变状态。
- 修改状态与修改标题、分类、截止时间等信息分开调用 `update`；传入 `--new-status` 时其他修改字段不会生效。
- 截止时间使用带时区的 RFC 3339 格式，例如 `2026-10-20T18:00:00+08:00`；用户未提供时区且无法从上下文确定时，先确认。
- 根据命令实际结果回复，不把生成的计划或总结当作已写入的工作记录。

## 常用命令

```bash
chronicle list --json
chronicle create "任务标题" -c "分类" --json
chronicle get <task_id> --json
chronicle log <task_id> "实际完成的工作"
chronicle update <task_id> --new-status done
chronicle summary --json
chronicle weekly-summary --json
```

其他参数按需查看 `chronicle <command> --help`。

## 计划与总结

生成每日计划时，读取活跃任务，结合截止时间、阻塞状态和用户优先级组织建议；缺少信息时不要推断任务已完成或自行设置截止时间。

生成工作总结时，使用 `summary [YYYY-MM-DD] --json` 或 `weekly-summary [YYYY-MM-DD] --json` 的实际记录。没有记录就说明暂无记录。
输出格式随用户需求调整；只有用户要求写入文件或外部应用时，才按指定位置保存。
