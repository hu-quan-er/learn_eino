# 第 31 课：自己实现一个最小 Agent

这一课只讲一件事：不用 `ChatModelAgent`，自己实现一个最小的 `adk.Agent`。

学完这一课，你要记住：

```text
Name / Description / Run -> Runner -> AgentEvent
```

## 1. 这一课解决什么问题

前面你已经会“使用” agent 了。

这一课开始，要把视角切到：

- agent 到底最少要实现什么
- `Runner` 到底在消费什么
- 不接模型时，agent 还能不能独立工作

答案是：可以。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson31-custom-agent/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson31-custom-agent
```

## 4. 这节课最关键的代码

### 第一步：实现 `Agent` 接口

```go
type customAgent struct {
    name string
}
```

最少只需要三件事：

- `Name(ctx)`
- `Description(ctx)`
- `Run(ctx, input, opts...)`

### 第二步：在 `Run(...)` 里自己产出事件

```go
event := adk.EventFromMessage(
    schema.AssistantMessage(content, nil),
    nil,
    schema.Assistant,
    "",
)
event.AgentName = a.name
```

这里的重点不是模型，而是：

- agent 的输出本质上是 `AgentEvent`
- 最常见的输出是 `MessageOutput`

### 第三步：交给 `Runner` 统一执行

```go
runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent})
iter := runner.Query(ctx, "请解释自定义 Agent 的最小实现")
```

这说明：

- 自定义 agent 和内置 agent 的顶层执行入口是一样的

## 5. 本课真正要记住的事

1. `Agent` 的最小实现核心就是 `Run(...)`
2. `Run(...)` 返回的是 `*adk.AsyncIterator[*adk.AgentEvent]`
3. `Runner` 不关心你是不是模型 agent，它只关心你有没有按协议产出事件
4. `adk.EventFromMessage(...)` 是构造消息事件的最简单方法
5. 先把本地 mock agent 写明白，再去接模型，学习曲线会平很多

## 6. 为什么这一课重要

因为后面你讲 `interrupt`、`resume`、`streaming`、`agent tool`，本质上都还是在扩展这一个接口模型。

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
