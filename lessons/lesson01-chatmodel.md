# 第 01 课：跑通第一个 ChatModel

这一课只做一件事：让你从零跑通一次 Eino 的模型调用。

学完这一课，你只需要记住 Eino 的最小调用路径：

```text
创建 ChatModel -> 准备 Message 列表 -> 调用 Generate -> 读取返回的 Message
```

## 1. 先理解 Eino 是什么

先不要把 Eino 想复杂。对入门阶段来说，你可以把它理解成一套 Go 里的 AI 应用积木体系：

- `components`：基础组件，比如 `ChatModel`、`Tool`、`Retriever`
- `schema`：统一的数据结构，比如消息、文档、工具参数
- `adk`：更高层的 Agent 能力
- `compose`：把多个组件串成工作流

第一课我们只碰 `ChatModel`，因为它是所有后续能力的入口。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson01-chatmodel/main.go`

这个 demo 用的是 Eino 官方生态里的 OpenAI ChatModel 实现，但它同样适用于很多 OpenAI 兼容接口。

## 3. 运行前准备

在当前目录执行：

```bash
cd /Users/seaya/huquan.95/project-4/eino
cp .env.example .env
```

然后准备环境变量。你可以手动 `export`，也可以自己用喜欢的方式加载 `.env`。

最少需要：

```bash
export OPENAI_API_KEY="你的密钥"
export OPENAI_MODEL="gpt-4o-mini"
export OPENAI_BASE_URL="你的兼容接口地址"
```

说明：

- `OPENAI_API_KEY`：模型服务密钥
- `OPENAI_MODEL`：模型名
- `OPENAI_BASE_URL`：可选；如果你不是直连 OpenAI 官方，大多数时候要填

## 4. 如何运行

```bash
go run ./cmd/lesson01-chatmodel
```

如果配置正常，你会看到类似输出：

```text
assistant:
Eino 是一个用 Go 构建大模型应用、工作流和 Agent 的开发框架。
```

## 5. 代码怎么读

### 第一步：初始化模型

```go
chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
    APIKey:  mustEnv("OPENAI_API_KEY"),
    Model:   mustEnv("OPENAI_MODEL"),
    BaseURL: os.Getenv("OPENAI_BASE_URL"),
})
```

这里做的是“把模型接进 Eino”。

你可以先把 `openai.NewChatModel(...)` 理解成：

- 输入一组模型配置
- 返回一个符合 Eino `ChatModel` 抽象的对象

后面不管你接 OpenAI、兼容 OpenAI 的服务、还是别的实现，只要它们遵守同一套接口，你的上层代码就更容易保持稳定。

### 第二步：准备消息

```go
messages := []*schema.Message{
    {
        Role:    schema.System,
        Content: "你是一个讲解 Eino 的入门老师，回答要短、准、清晰。",
    },
    {
        Role:    schema.User,
        Content: "请用一句话解释什么是 Eino。",
    },
}
```

这是 Eino 里非常核心的一层：`schema.Message`。

你现在只需要先掌握两种角色：

- `schema.System`：系统设定
- `schema.User`：用户输入

后面你还会见到 `schema.Assistant`、工具调用相关消息、多模态消息等。

### 第三步：调用 Generate

```go
resp, err := chatModel.Generate(ctx, messages, model.WithTemperature(0.2))
```

这里是本课最关键的一行。

它的意思是：

- 把消息列表交给模型
- 让模型生成一个完整回答
- 通过 `model.WithTemperature(0.2)` 传入本次调用参数

你可以先这样理解：

- `NewChatModel(...)` 负责创建“模型实例”
- `Generate(...)` 负责发起“一次具体调用”
- `model.WithTemperature(...)` 这类参数负责调整“这次调用”的行为

### 第四步：读取返回值

```go
fmt.Println(resp.Content)
```

返回值本质上还是一个 `schema.Message`。也就是说，Eino 的输入和输出都尽量统一在同一套 `schema` 结构里，这会让后续组合组件时更自然。

## 6. 本课你真正需要记住的事

只记住这四点：

1. Eino 先提供统一抽象，再接不同模型实现
2. `ChatModel` 是最基础的模型组件
3. `schema.Message` 是最核心的数据结构之一
4. `Generate` 用于一次性拿完整结果

## 7. 你现在可以自己做的练习

建议你立刻改三处，再重新运行：

1. 改 `System` 提示词，让回答风格变得更正式
2. 改 `User` 问题，比如问 “Eino 和 LangChain 有什么关系”
3. 改 `temperature`，观察回答稳定性变化

## 8. 下一课会讲什么

下一课最自然的延续是两种方向：

- `Stream`：为什么流式输出和 `Generate` 不一样
- `Prompt`：如何把写死的消息变成可复用模板

如果你要继续，我建议先讲 `Stream`，因为它仍然围绕 `ChatModel`，认知跨度最小。

## 9. 官方资料

本课内容参考的是 Eino 官方一手资料：

- Quick Start: https://www.cloudwego.io/docs/eino/quick_start/
- ChatModel Guide: https://www.cloudwego.io/docs/eino/core_modules/components/chat_model_guide/
- OpenAI ChatModel: https://www.cloudwego.io/docs/eino/ecosystem_integration/chat_model/chat_model_openai/
- Eino 仓库: https://github.com/cloudwego/eino
