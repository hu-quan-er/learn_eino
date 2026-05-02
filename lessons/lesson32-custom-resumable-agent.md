# 第 32 课：自己实现一个 ResumableAgent

这一课只讲一件事：自己把 `interrupt -> checkpoint -> resume` 做到自定义 agent 里。

学完这一课，你要记住：

```text
Run -> StatefulInterrupt -> Runner checkpoint -> Resume(ctx, info)
```

## 1. 这一课解决什么问题

如果你只是“用”内置 agent，中断恢复看起来像框架魔法。

但真正做复杂系统时，你必须清楚：

- 中断点是谁发出来的
- 恢复时状态是谁接回去的
- `ResumeInfo` 里到底有什么

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson32-custom-resumable-agent/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson32-custom-resumable-agent
```

你会看到两段输出：

1. 第一次运行先中断并打印 interrupt root cause
2. 第二次从 checkpoint 恢复并完成审批

## 4. 这节课最关键的代码

### 第一步：让 agent 同时实现 `Run` 和 `Resume`

```go
type resumableAgent struct {
    name string
}
```

只实现 `Run` 还不够。

如果你想让 `Runner.Resume...` 能继续执行，这个 agent 本身就必须实现 `Resume(...)`。

### 第二步：第一次运行里主动发出 `StatefulInterrupt`

```go
interruptEvent := adk.StatefulInterrupt(ctx, "需要人工审批", ApprovalState{Request: request})
```

这里有两个层次：

- `info`：给外部看的中断原因
- `state`：恢复时自己继续用的内部状态

### 第三步：状态类型要注册

```go
func init() {
    schema.Register[ApprovalState]()
}
```

这是 checkpoint 场景里非常容易漏掉的一步。

因为中断状态会经过 gob 序列化。

### 第四步：恢复时从 `ResumeInfo` 取数据

```go
state, ok := approvalStateFromAny(info.InterruptState)
decision, _ := info.ResumeData.(string)
```

`ResumeInfo` 里最重要的是：

- `WasInterrupted`
- `InterruptState`
- `IsResumeTarget`
- `ResumeData`

### 第五步：第一次运行必须带 `WithCheckPointID`

```go
iter := runner.Query(ctx, "发布第 32 课", adk.WithCheckPointID(checkPointID))
```

如果没有 checkpoint ID，`Runner` 没有地方存恢复点。

## 5. 本课真正要记住的事

1. 想恢复自定义 agent，就必须实现 `ResumableAgent`
2. `StatefulInterrupt` 适合“恢复时还要拿回内部状态”的场景
3. 自定义状态要先 `schema.Register[...]()`
4. `ResumeWithParams(...)` 可以定点恢复并携带恢复数据
5. `ResumeInfo` 才是恢复阶段真正该读的上下文

## 6. 这节课和第 14 课的区别

第 14 课重点是：

- “使用”已有 agent 和 tool 跑通恢复闭环

第 32 课重点是：

- “实现”你自己的 resumable agent

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
