# 第 06 课：学会 ToolsNode 执行工具调用

这一课只讲一件事：理解 Eino 里“工具定义”和“工具执行”是分开的。

学完这一课，你要记住：

```text
assistant tool call -> ToolsNode -> []*schema.Message(tool results)
```

## 1. 这一课解决什么问题

第 04 课你已经会定义 Tool 了，但那时候还是你自己直接调：

- `tool.Info()`
- `tool.InvokableRun(...)`

这还不是 Eino 在运行时真正执行工具的方式。

真正负责执行工具调用的是 `ToolsNode`。

它的职责是：

1. 接收一条带 `ToolCalls` 的 assistant message
2. 找到对应工具
3. 执行工具
4. 产出标准 `tool` message

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson06-tools-node/main.go`

这节课是纯本地 demo，不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson06-tools-node
```

你会看到类似输出：

```text
tool messages:
1. role=tool tool_name=get_weather tool_call_id=call_weather_1
   content={"summary":"Shanghai is sunny today"}
```

## 4. 这节课最关键的代码

### 第一步：构造一个工具调用消息

```go
assistantMessage := schema.AssistantMessage("", []schema.ToolCall{
    {
        ID:   "call_weather_1",
        Type: "function",
        Function: schema.FunctionCall{
            Name:      "get_weather",
            Arguments: `{"city":"Shanghai"}`,
        },
    },
})
```

这条消息不是普通回答，而是“模型决定要调工具”的消息。

关键点在于：

- `Role` 还是 `assistant`
- 但消息里带了 `ToolCalls`

这就是模型请求工具执行的标准形态。

### 第二步：创建 ToolsNode

```go
toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
    Tools: []any{weatherTool},
})
```

这里本质上是在说：

- 当前这个执行节点可以调用哪些工具

### 第三步：执行工具调用

```go
toolMessages, err := toolsNode.Invoke(ctx, assistantMessage)
```

这是本课最关键的一行。

`Invoke(...)` 会把 assistant message 里的每个 tool call 执行掉，然后返回：

- `[]*schema.Message`

这些返回值的 `Role` 会变成：

- `schema.Tool`

## 5. 这节课真正要记住的事

只记这四点：

1. Tool 是能力定义
2. ToolsNode 是能力执行器
3. ToolsNode 的输入通常是带 `ToolCalls` 的 assistant message
4. ToolsNode 的输出是标准 tool messages

## 6. 下一课讲什么

下一课把模型和 ToolsNode 串起来，跑通完整顺序：

- 模型先产出 tool call
- ToolsNode 执行工具
- 模型再基于工具结果给最终答案

## 7. 官方资料

- Tools Guide: https://www.cloudwego.io/docs/eino/core_modules/components/tools_node_guide/
- Eino 仓库: https://github.com/cloudwego/eino
