# 第 07 课：跑通模型调用工具的最小闭环

这一课只讲一件事：把模型和工具真正连起来。

学完这一课，你要记住：

```text
model -> assistant tool call -> ToolsNode -> tool message -> model final answer
```

## 1. 这一课解决什么问题

到第 06 课为止，你已经分别见过：

- Tool 怎么定义
- ToolsNode 怎么执行工具调用

但这两个东西还没有和真实模型连起来。

所以这一课的目标很明确：

1. 让模型先决定要不要调工具
2. 执行工具
3. 把工具结果再喂回模型
4. 拿到最终自然语言回答

这就是很多 Agent 系统最基础的运行循环。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson07-model-with-tools/main.go`

这节课需要模型密钥，因为第一步和最后一步都要真实调用模型。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
export OPENAI_API_KEY="你的密钥"
export OPENAI_MODEL="gpt-4o-mini"
export OPENAI_BASE_URL="你的兼容接口地址"
go run ./cmd/lesson07-model-with-tools
```

你会先看到模型产出的 tool call，再看到最后的自然语言答案。

## 4. 这节课最关键的代码

### 第一步：把工具 schema 告诉模型

```go
toolInfo, err := weatherTool.Info(ctx)
firstResp, err := chatModel.Generate(
    ctx,
    history,
    model.WithTools([]*schema.ToolInfo{toolInfo}),
    model.WithToolChoice(schema.ToolChoiceForced),
)
```

这里有两件事要分清：

1. `model.WithTools(...)` 是把工具描述传给模型
2. `model.WithToolChoice(...)` 是控制模型是否必须调工具

这节课用 `ToolChoiceForced`，目的是让 demo 更稳定。

### 第二步：检查模型产出的 ToolCalls

```go
for _, toolCall := range firstResp.ToolCalls {
    ...
}
```

如果模型决定调工具，它返回的不是最终答案，而是一条 assistant message，里面带：

- `ToolCalls`

### 第三步：执行工具

```go
toolMessages, err := toolsNode.Invoke(ctx, firstResp)
```

这一步就是第 06 课学过的内容。

### 第四步：把工具结果接回模型

```go
history = append(history, firstResp)
history = append(history, toolMessages...)
finalResp, err := chatModel.Generate(ctx, history)
```

这是这一课最关键的闭环。

你要记住：工具执行完以后，模型并不会自动知道结果，除非你把 tool messages 放回对话历史里。

## 5. 本课真正要记住的事

只记这五点：

1. 模型看到的是 `ToolInfo`
2. 模型返回工具请求时，消息里会带 `ToolCalls`
3. `ToolsNode` 负责执行这些调用
4. 执行结果要转成 tool messages
5. tool messages 必须再喂回模型，模型才能给最终答案

## 6. 下一课讲什么

下一课先从 `Workflow` 开始，继续进入编排层。

因为到这一步你已经有了组件和组件之间的闭环，下一步应该学会更通用的图式编排。

## 7. 官方资料

- ChatModel Guide: https://www.cloudwego.io/docs/eino/core_modules/components/chat_model_guide/
- Tools Guide: https://www.cloudwego.io/docs/eino/core_modules/components/tools_node_guide/
- Eino 仓库: https://github.com/cloudwego/eino
