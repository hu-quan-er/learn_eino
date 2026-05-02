# 第 10 课：跑通第一个 ChatModelAgent

这一课只讲一件事：让你第一次跑通 Eino ADK 的 agent 执行入口。

学完这一课，你要记住：

```text
ChatModelAgent -> Runner -> AgentEvent
```

## 1. 这一课解决什么问题

前面几课你学到的还是组件级能力：

- Prompt
- Model
- Tool
- ToolsNode
- Chain
- Workflow

但到了应用层，你不一定想自己手写整套“模型循环 + 工具执行 + 事件处理”。

这时候 Eino ADK 提供了更高层入口：

- `ChatModelAgent`
- `Runner`

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson10-chatmodel-agent/main.go`

这节课需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
export OPENAI_API_KEY="你的密钥"
export OPENAI_MODEL="gpt-4o-mini"
export OPENAI_BASE_URL="你的兼容接口地址"
go run ./cmd/lesson10-chatmodel-agent
```

你会看到一串事件输出。通常会先看到 tool 相关事件，再看到最终 assistant 回答。

## 4. 这节课最关键的代码

### 第一步：创建 Agent

```go
agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Name:        "weather_agent",
    Description: "answer weather questions with a weather tool",
    Instruction: "你是一个天气助教。必要时调用 get_weather 工具。拿到工具结果后，用中文给出简短答案。",
    Model:       chatModel,
    ToolsConfig: adk.ToolsConfig{
        ToolsNodeConfig: compose.ToolsNodeConfig{
            Tools: []any{weatherTool},
        },
    },
})
```

这里你要先记住三件事：

1. `Instruction` 是 agent 的系统指令
2. `Model` 是底层聊天模型
3. `ToolsConfig` 是 agent 能调用哪些工具

### 第二步：创建 Runner

```go
runner := adk.NewRunner(ctx, adk.RunnerConfig{
    Agent: agent,
})
```

`Runner` 是 ADK 的执行入口。

你可以先把它理解成：

- agent 的运行器

### 第三步：发起 Query

```go
iter := runner.Query(ctx, "Shanghai 的天气怎么样？")
```

这一行会启动整条 agent 执行流程。

### 第四步：读取 AgentEvent

```go
for {
    event, ok := iter.Next()
    ...
}
```

这是这一课最关键的动作。

ADK 不只是给你一个最终字符串，而是给你事件流：

- 有的事件表示 tool 输出
- 有的事件表示 assistant 输出
- 还有可能带错误或动作信息

### 第五步：拿到最终消息

```go
message, err := event.Output.MessageOutput.GetMessage()
```

`GetMessage()` 的好处是：

- 不管底层是普通 message 还是 streaming message
- 都能统一拿到最终消息

## 5. 本课真正要记住的事

只记这五点：

1. `ChatModelAgent` 是更高层的 agent 封装
2. `Runner` 是执行入口
3. `Query(...)` 是最简单的发起方式
4. 执行结果不是单个值，而是一串 `AgentEvent`
5. `MessageOutput.GetMessage()` 可以统一读取最终消息

## 6. 到第 10 课为止，你已经掌握了什么

现在你已经具备了一整条 Eino 入门主线：

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

这时你再继续往后学：

- 中断与恢复
- 多 agent
- 更复杂的 workflow

会容易很多。

## 7. 官方资料

- Eino ADK 文档入口：https://www.cloudwego.io/docs/eino/
- Eino 仓库: https://github.com/cloudwego/eino
