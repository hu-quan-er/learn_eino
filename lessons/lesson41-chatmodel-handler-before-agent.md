# 第 41 课：学会 ChatModelAgentMiddleware 的 BeforeAgent

这一课只讲一件事：在每次运行开始前，动态修改 instruction 和 tools。

学完这一课，你要记住：

```text
BeforeAgent -> mutate Instruction / Tools / ReturnDirectly
```

## 1. 这一课解决什么问题

`AgentMiddleware` 适合静态扩展。

如果你想做的是：

- 运行时动态加 tool
- 运行时改 instruction
- 不同请求走不同 agent 配置

那就该看 `ChatModelAgentMiddleware.BeforeAgent(...)`。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson41-chatmodel-handler-before-agent/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson41-chatmodel-handler-before-agent
```

## 4. 这节课最关键的代码

### 自定义 handler

```go
type policyInjectionHandler struct {
    *adk.BaseChatModelAgentMiddleware
    tool tool.BaseTool
}
```

### 在 `BeforeAgent(...)` 里动态改运行上下文

```go
runCtx.Instruction += "\n\n运行时规则：必须调用 lookup_policy 工具后再回答。"
runCtx.Tools = append(runCtx.Tools, h.tool)
```

这个修改不是全局配置，而是当前 run 的运行态。

## 5. 本课真正要记住的事

1. `BeforeAgent(...)` 是 ChatModelAgent 运行前最重要的动态配置入口
2. 可以改 instruction
3. 可以改 tools
4. 可以按请求动态注入能力
5. 这是很多生产 middleware 的起点

## 6. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
