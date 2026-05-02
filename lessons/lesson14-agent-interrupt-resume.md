# 第 14 课：学会 Agent 中断与恢复

这一课只讲一件事：跑通一次真正的 `interrupt -> checkpoint -> resume` 闭环。

学完这一课，你要记住：

```text
tool.StatefulInterrupt -> Runner checkpoint -> ResumeWithParams
```

## 1. 这一课解决什么问题

很多 agent 任务不是模型自己就能跑完的。

比如：

- 需要人工审批
- 需要外部系统确认
- 需要用户补参数

这时候你不能只返回错误，也不能硬编码阻塞等待。

正确做法是：

1. 先中断
2. 保存 checkpoint
3. 等外部信息回来后恢复

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson14-agent-interrupt-resume/main.go`

这节课不需要模型密钥。demo 里用了一个本地脚本化 model 和一个会中断的审批 tool。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson14-agent-interrupt-resume
```

你会看到两段输出：

1. 第一次运行先触发 interrupt
2. 第二次通过 `approved` 恢复执行

## 4. 这节课最关键的代码

### 第一步：tool 主动中断

```go
return "", tool.StatefulInterrupt(
    ctx,
    fmt.Sprintf("需要审批：%s", input.Action),
    "pending:"+input.Action,
)
```

这里有两个值：

- `info`：给外部看的中断原因
- `state`：给恢复时自己用的内部状态

### 第二步：Runner 要挂 checkpoint store

```go
runner := adk.NewRunner(ctx, adk.RunnerConfig{
    Agent:           agent,
    CheckPointStore: newInMemoryStore(),
})
```

如果没有 store，就没有恢复能力。

### 第三步：从 interrupt 事件里拿恢复目标 ID

```go
for _, interruptCtx := range event.Action.Interrupted.InterruptContexts {
    if interruptCtx.IsRootCause {
        rootCauseID = interruptCtx.ID
    }
}
```

这个 `ID` 很关键。

它不是普通日志字段，而是恢复时要精确命中的 target。

### 第四步：带参数恢复

```go
iter, err = runner.ResumeWithParams(ctx, checkpointID, &adk.ResumeParams{
    Targets: map[string]any{
        rootCauseID: "approved",
    },
})
```

这表示：

- 恢复这个具体 interrupt 点
- 恢复数据是 `"approved"`

### 第五步：在 tool 里读取恢复上下文

```go
isResumeTarget, hasData, decision := tool.GetResumeContext[string](ctx)
```

如果当前 tool 就是这次 resume 的目标，它就能拿到：

- 自己是不是被点名恢复
- 有没有带数据
- 恢复数据是什么

## 5. 本课真正要记住的事

1. 中断入口通常发生在 tool 或 node 内部
2. `StatefulInterrupt` 会同时保存中断状态
3. `Runner` 要配 `CheckPointStore`
4. 恢复时要用 root cause interrupt 的 `ID`
5. `ResumeWithParams(...)` 适合“定点恢复 + 带数据恢复”

## 6. 和第 15 课的区别

第 14 课重点是：

- ADK agent 层的中断恢复
- 精确命中某个 interrupt target

第 15 课重点是：

- compose workflow 层的 checkpoint 恢复
- 更像流程编排层的暂停与继续

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
