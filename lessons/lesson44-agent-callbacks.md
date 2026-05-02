# 第 44 课：学会用 WithCallbacks 观察 Agent 运行

这一课只讲一件事：用 `WithCallbacks(...)` 在 agent 运行开始和结束时拿到结构化信息。

学完这一课，你要记住：

```text
WithCallbacks -> OnStart -> OnEnd -> ConvAgentCallbackInput / Output
```

## 1. 这一课解决什么问题

当你要做：

- 运行日志
- trace
- 统计
- 调试

你不能只靠业务代码里手动 `fmt.Println`。

更稳的入口是 callback。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson44-agent-callbacks/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson44-agent-callbacks
```

## 4. 这节课最关键的代码

### 注册 callback handler

```go
handler := callbacks.NewHandlerBuilder().
    OnStartFn(...).
    OnEndFn(...).
    Build()
```

### 启动运行时挂上 callback

```go
runner.Query(ctx, "lesson44 callback demo", adk.WithCallbacks(handler))
```

### 转换输入输出类型

```go
adk.ConvAgentCallbackInput(input)
adk.ConvAgentCallbackOutput(output)
```

因为 callback 的签名是通用接口，所以 agent 侧要做一次类型转换。

## 5. 本课真正要记住的事

1. `WithCallbacks(...)` 是 agent 运行级观测入口
2. `OnStart` 适合拿 run info 和输入
3. `OnEnd` 适合拿事件流副本
4. callback 拿到的事件流应该异步消费，避免阻塞主执行
5. 生产里做 tracing / logging / metrics，callback 很重要

## 6. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
