# log-daily-work

获取今日工作记录，写入 Obsidian Daily Log。

## CLI 命令

```powershell
chronicle summary [date] --json
```

## 执行步骤

### 1. 获取今日记录

```powershell
chronicle summary --json
```

### 2. 处理

- 无记录：提示用户
- 有记录：解析 JSON，浓缩成一句话摘要

### 3. 写入文件

目标文件：`Daily {YYYYMM}.md`

格式：
```
| Date | Work | Hour | Comment |
| ---- | ---- | ---- | ------- |
| 20260313 | **任务标题**：一句话摘要 | | |
```

## 使用触发

- "log daily"
- "记录今天的工作"
- "update daily"
- "更新日报"
- "log work"
