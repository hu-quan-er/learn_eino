# 第 18 课：把 Agent 包装成 Tool

这一课只讲一件事：把一个 agent 变成上层可调用的 tool。

学完这一课，你要记住：

```text
Agent -> NewAgentTool -> ToolsNode -> ToolMessage
```

## 1. 这一课解决什么问题

当你有一个已经写好的 agent 时，很多时候你希望：

- 不直接把它当顶层入口跑
- 而是把它作为某个更大系统里的一个“能力模块”

这就是 `AgentTool` 的用途。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson18-agent-tool/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson18-agent-tool
```

## 4. 这节课最关键的代码

### 第一步：先有一个 Agent

```go
type FAQAgent struct{}
```

本课为了讲清楚概念，用了一个最小自定义 agent。

### 第二步：把 agent 包成 tool

```go
agentTool := adk.NewAgentTool(ctx, &FAQAgent{})
```

从这一步开始，这个 agent 就能像普通 tool 一样被上层调用。

### 第三步：把它塞进 ToolsNode

```go
toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
    Tools: []tool.BaseTool{agentTool},
})
```

这说明：

- `AgentTool` 最终对外暴露的仍然是 `tool.BaseTool`

### 第四步：像 tool call 一样执行

```go
toolMessages, err := toolsNode.Invoke(ctx, assistantMessage)
```

最终返回的仍然是标准 `tool message`。

也就是说，上层根本不需要知道这个 tool 背后其实是 agent。

## 5. 本课真正要记住的事

1. `NewAgentTool(...)` 可以把 agent 变成 tool
2. 包装后的 agent 可以进入 `ToolsNode`
3. 上层调用方式和普通 tool 一样
4. `AgentTool` 适合把已有 agent 复用到更大系统里
5. 这是多 agent 组合的重要基础

## 6. 为什么这节课重要

因为从这一课开始，你不再只是：

- 模型调用普通函数型 tool

而是可以做到：

- 一个 agent 调另一个 agent

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
