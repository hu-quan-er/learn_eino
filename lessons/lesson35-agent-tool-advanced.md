# 第 35 课：学会 AgentTool 的高级输入模式

这一课只讲一件事：把 `AgentTool` 的三种关键输入模式讲清楚。

学完这一课，你要记住：

```text
NewAgentTool -> default request / custom schema / full chat history
```

## 1. 这一课解决什么问题

很多人会把 `AgentTool` 只理解成：

- “把 agent 包成 tool”

这还不够。

真正会影响你系统设计的是：

- tool 参数怎么进 agent
- 是只给一个 `request`
- 还是给原始 JSON
- 还是把整个父级聊天历史一起带进去

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson35-agent-tool-advanced/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson35-agent-tool-advanced
```

你会看到三段输出：

1. default schema
2. custom schema
3. full chat history

## 4. 这节课最关键的代码

### 第一种：默认输入模式

```go
agentTool := adk.NewAgentTool(ctx, agent)
result, err := agentTool.(tool.InvokableTool).InvokableRun(ctx, `{"request":"请解释默认 request 字段"}`)
```

默认情况下：

- tool schema 是 `{"request": string}`
- 子 agent 收到的是一个普通 `user message`
- message content 就是 `request` 的值

### 第二种：自定义输入 schema

```go
agentTool := adk.NewAgentTool(ctx, agent, adk.WithAgentInputSchema(customSchema))
```

这一种要记住一个关键行为：

- 子 agent 收到的不是拆开的字段
- 而是“整段原始 JSON 字符串”

也就是：

```json
{"topic":"agent tool","level":"advanced"}
```

会原样进到子 agent 的 message content 里。

### 第三种：完整聊天历史作为输入

```go
childTool := adk.NewAgentTool(ctx, childAgent, adk.WithFullChatHistoryAsInput())
```

这种模式下：

- tool 当前参数不再是重点
- 框架会从父级上下文中提取聊天历史
- 再追加 transfer messages
- 然后把这一整段 history 作为子 agent 输入

这正是“一个 agent 把任务转交给另一个 agent”时最有用的模式。

## 5. 本课 demo 里一个值得注意的点

为了让 `WithFullChatHistoryAsInput()` 在本地 demo 里可观察，本课用了一个小 graph，把 `adk.State.Messages` 放进 graph local state。

这在当前版本能工作，但你也要知道：

- `adk.State` 源码里已经标了未来会逐步收口
- 所以生产代码里更稳的方向还是跟着后续官方 handler / middleware 方式演进

## 6. 本课真正要记住的事

1. 默认模式下，子 agent 收到的是 `request` 字段值
2. 自定义 schema 模式下，子 agent 收到的是原始 JSON 字符串
3. `WithFullChatHistoryAsInput()` 适合真正的 agent-to-agent 任务转交
4. `AgentTool` 不只是“包一层”，它决定了子 agent 的输入语义
5. 设计多 agent 系统时，输入模式选错，后面 prompt 和状态传递都会越来越乱

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
