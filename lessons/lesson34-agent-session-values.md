# 第 34 课：用 Session Values 共享 Agent 运行态

这一课只讲一件事：理解 ADK `session values` 在多 agent 运行中的作用。

学完这一课，你要记住：

```text
WithSessionValues -> AddSessionValue -> GetSessionValue / GetSessionValues
```

## 1. 这一课解决什么问题

很多人第一次做多 agent，会把共享数据全塞回 message。

这不是最好的做法。

因为有些数据：

- 不是给模型看的
- 不是对话历史
- 只是这次运行过程里的内部上下文

这种数据更适合放进 `session values`。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson34-agent-session-values/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson34-agent-session-values
```

## 4. 这节课最关键的代码

### 第一步：启动时注入 session values

```go
iter := runner.Query(ctx, "准备一节 ADK 教程", adk.WithSessionValues(map[string]any{
    "tenant": "tutorial-team",
}))
```

这表示：

- 一次运行开始前，就可以先给 session 塞公共上下文

### 第二步：前一个 agent 写入值

```go
adk.AddSessionValue(ctx, "plan", plan)
adk.AddSessionValue(ctx, "owner", a.name)
```

这里写进去的是：

- 当前运行里后续 agent 可读的内部数据

### 第三步：后一个 agent 读取值

```go
plan, _ := adk.GetSessionValue(ctx, "plan")
keys := sessionKeys(adk.GetSessionValues(ctx))
```

这说明：

- 同一次 run 里的后续 agent，可以直接读到前面写入的 session

## 5. 本课真正要记住的事

1. `session values` 是运行时上下文，不是消息历史
2. 它适合存内部状态、路由信息、计划、租户信息这类数据
3. `WithSessionValues(...)` 适合 run 入口注入公共上下文
4. `AddSessionValue(...)` 适合 agent 运行中补充共享数据
5. 多 agent 协作时，session 往往比硬塞 message 更干净

## 6. 它和 Message History 的区别

`message history`：

- 主要给模型看
- 会影响上下文窗口

`session values`：

- 主要给程序逻辑看
- 不一定需要进入模型上下文

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
