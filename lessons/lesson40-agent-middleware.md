# 第 40 课：学会 AgentMiddleware 的基础钩子

这一课只讲一件事：先把 `AgentMiddleware` 这套老而直接的扩展方式看明白。

学完这一课，你要记住：

```text
AdditionalInstruction / BeforeChatModel / AfterChatModel
```

## 1. 这一课解决什么问题

当你只需要做一些简单扩展时，比如：

- 补一句 instruction
- 在模型前后轻度改 state

`AgentMiddleware` 已经够用。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson40-agent-middleware/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson40-agent-middleware
```

## 4. 这节课最关键的代码

### `AdditionalInstruction`

```go
AdditionalInstruction: "额外规则：回答里必须显式出现 middleware。",
```

它会追加到 agent instruction 里。

### `BeforeChatModel`

```go
state.Messages = append(state.Messages, schema.UserMessage("before:middleware"))
```

它会在模型调用前改写当前 state。

### `AfterChatModel`

```go
last.Content += " | after:middleware"
```

它会在模型返回后继续改 state。

## 5. 本课真正要记住的事

1. `AgentMiddleware` 适合简单、静态扩展
2. 它能在模型调用前后改 `ChatModelAgentState`
3. `AdditionalInstruction` 是最轻量的插桩点
4. 如果你需要更灵活的上下文传播和 wrapper，后面要看 `ChatModelAgentMiddleware`
5. `AgentMiddleware` 现在仍然有用，但更偏基础能力

## 6. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
