# 第 15 课：给 Workflow 加上 Checkpoint

这一课只讲一件事：让 workflow 在跑到关键节点前先暂停，随后从 checkpoint 继续执行。

学完这一课，你要记住：

```text
Workflow + CheckPointStore + CheckPointID -> interrupt -> resume
```

## 1. 这一课解决什么问题

如果你做的是流程编排，不一定总需要 agent 级别的 `ResumeWithParams`。

很多时候你只想要：

- 跑到某个节点前先停下
- 等外部确认后继续

这就是 workflow checkpoint 的典型场景。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson15-workflow-checkpoint/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson15-workflow-checkpoint
```

你会看到：

1. 第一次运行在 `send_notification` 前中断
2. 第二次用同一个 checkpoint ID 继续执行

## 4. 这节课最关键的代码

### 第一步：编译 workflow 时挂上 checkpoint 能力

```go
runner, err := workflow.Compile(ctx,
    compose.WithCheckPointStore(store),
    compose.WithInterruptBeforeNodes([]string{"send_notification"}),
)
```

这里同时做了两件事：

1. 设置 checkpoint store
2. 指定在 `send_notification` 前中断

### 第二步：第一次运行时带上 checkpoint ID

```go
_, err = runner.Invoke(ctx, map[string]any{
    "topic": "今晚 8 点发布 Eino 第 15 课",
}, compose.WithCheckPointID(checkpointID))
```

如果命中中断点，这里不会正常返回结果，而是返回 interrupt error。

### 第三步：读取 interrupt 信息

```go
info, ok := compose.ExtractInterruptInfo(err)
```

你可以从这里看到：

- 哪些节点前被中断
- 当前 interrupt address

### 第四步：继续执行

```go
output, err := runner.Invoke(ctx, nil, compose.WithCheckPointID(checkpointID))
```

这一行是本课最重要的地方。

这里没有重新传业务输入，而是直接靠 checkpoint 恢复。

也就是说：

- 已经跑过的部分不需要重新计算
- workflow 会从上次暂停的位置继续

## 5. 本课真正要记住的事

1. compose 层也支持 checkpoint
2. `WithCheckPointStore(...)` 决定状态存哪里
3. `WithCheckPointID(...)` 决定恢复哪次执行
4. `WithInterruptBeforeNodes(...)` 适合做流程级暂停点
5. 恢复时可以继续用同一个 runner

## 6. 到第 15 课为止，你已经掌握了什么

现在你已经把 Eino 的入门主线走到比较完整了：

1. ChatModel
2. Stream
3. Prompt
4. Tool
5. Chain
6. ToolsNode
7. Model + Tool 闭环
8. Workflow
9. MessageParser
10. ChatModelAgent
11. StreamableTool
12. Workflow Branch
13. Agent Stream
14. Agent Interrupt / Resume
15. Workflow Checkpoint

后面再往下，可以继续进：

- 多 agent 协作
- 更复杂的 graph / state
- 一个完整项目骨架

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
