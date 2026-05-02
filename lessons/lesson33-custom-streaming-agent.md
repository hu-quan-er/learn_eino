# 第 33 课：自己实现一个 Streaming Agent

这一课只讲一件事：让自定义 agent 直接产出流式 `MessageStream`。

学完这一课，你要记住：

```text
MessageVariant{IsStreaming:true} -> MessageStream -> Runner(EnableStreaming:true)
```

## 1. 这一课解决什么问题

前面你已经看过模型流式输出，也看过 agent 流式事件。

这一课要搞清楚的是：

- 如果 agent 不是 `ChatModelAgent`
- 而是你自己实现的
- 流式输出应该怎么拼

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson33-custom-streaming-agent/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson33-custom-streaming-agent
```

## 4. 这节课最关键的代码

### 第一步：先看 `input.EnableStreaming`

```go
if !input.EnableStreaming {
    ...
}
```

这说明：

- agent 是否走流式，不是你硬编码决定的
- 而是 `RunnerConfig.EnableStreaming` 传进来的

### 第二步：自己构造一个 `MessageStream`

```go
stream := schema.StreamReaderFromArray([]*schema.Message{
    schema.AssistantMessage("自定义 ", nil),
    schema.AssistantMessage("streaming ", nil),
    schema.AssistantMessage("agent", nil),
})
stream.SetAutomaticClose()
```

这里要记住：

- 流式 agent 的核心不是 `Generate`
- 而是 `MessageStream`

### 第三步：把 stream 放进 `MessageVariant`

```go
event := adk.EventFromMessage(nil, stream, schema.Assistant, "")
```

这一步会生成一个 `IsStreaming=true` 的消息事件。

### 第四步：消费端仍然是逐 chunk 读取

```go
chunk, err := output.MessageStream.Recv()
```

agent 的输出虽然不是模型直接来的，但消费方式和之前保持一致。

## 5. 本课真正要记住的事

1. 自定义 agent 完全可以自己产出 `MessageStream`
2. `RunnerConfig.EnableStreaming=true` 会传到 `AgentInput`
3. `SetAutomaticClose()` 是流式 message stream 的好习惯
4. 输出端看见的仍然是 `AgentEvent`
5. 流式 agent 和流式 model 在消费侧是统一心智模型

## 6. 为什么这一课重要

因为后面你做：

- agent 包 tool
- tool 再往上透出内部事件
- 多 agent 边跑边回显

本质上都依赖你理解 `AgentEvent` 里的 streaming message。

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
