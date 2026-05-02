# 第 13 课：跑通 Agent 的流式事件

这一课只讲一件事：让你真正看到 ADK `Runner` 在流式模式下返回的事件长什么样。

学完这一课，你要记住：

```text
Runner(EnableStreaming=true) -> AgentEvent -> MessageStream
```

## 1. 这一课解决什么问题

第 10 课你已经知道 agent 的执行结果是一串 `AgentEvent`。

但那时还是普通模式。

流式模式下要多理解一层：

- event 里不一定是完整 message
- 也可能是一个 `MessageStream`

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson13-agent-stream/main.go`

这节课不需要模型密钥。demo 里用了一个本地 mock model，专门演示流式事件。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson13-agent-stream
```

你会看到：

- 先打印流式 chunk
- 再打印拼回来的完整答案

## 4. 这节课最关键的代码

### 第一步：给 Runner 打开流式模式

```go
runner := adk.NewRunner(ctx, adk.RunnerConfig{
    Agent:           agent,
    EnableStreaming: true,
})
```

这一步是总开关。

### 第二步：判断当前事件是不是流式消息

```go
if output.IsStreaming {
    ...
}
```

这一步很重要。

因为 `AgentEvent` 可能同时承载两种结果：

- 普通 message
- streaming message

### 第三步：直接消费 `MessageStream`

```go
chunk, err := output.MessageStream.Recv()
```

你要把这一步和第 02 课联系起来理解：

- 模型流式输出时，最底层还是 `Recv()`

只不过现在这个 stream 被包进了 `AgentEvent` 里。

### 第四步：把多个 chunk 合并成最终消息

```go
finalMessage, err := schema.ConcatMessages(chunks)
```

这一步和普通流式模型一样，都是“边收边展示，最后再合并”。

## 5. 本课真正要记住的事

1. Agent 流式输出要在 `RunnerConfig` 里打开 `EnableStreaming`
2. 事件结果要先看 `MessageOutput.IsStreaming`
3. 如果是流式事件，要读 `MessageOutput.MessageStream`
4. 读完 chunk 后可以自己拼成完整消息
5. `AgentEvent` 是 agent 层的统一事件协议

## 6. 和第 02 课的区别

第 02 课是：

- 你直接面对 `ChatModel.Stream(...)`

第 13 课是：

- 你面对的是 `AgentEvent`
- 流被包在 event 里面

所以这节课本质是在学：

- agent 层怎么承接底层流式模型

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
