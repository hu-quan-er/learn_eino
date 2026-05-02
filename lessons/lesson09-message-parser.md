# 第 09 课：把 Message 解析成结构体

这一课只讲一件事：把一条 `schema.Message` 里的 JSON 内容解析成 Go 结构体。

学完这一课，你要记住：

```text
message -> MessageJSONParser -> Go struct
```

## 1. 这一课解决什么问题

前面的模型调用大多数还是自然语言输出：

- 看起来方便
- 但程序不一定好处理

很多时候你真正想要的是：

- 一个结构体
- 一个数组
- 一个固定字段的对象

这就是 `MessageParser` 的作用。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson09-message-parser/main.go`

这节课用本地构造的 assistant message 做解析，不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson09-message-parser
```

你会看到类似输出：

```text
parsed struct: {Topic:ToolsNode KeyPoints:[executes tool calls returns tool messages] Difficulty:2}
```

## 4. 这节课最关键的代码

### 第一步：定义目标结构体

```go
type LessonSummary struct {
    Topic      string   `json:"topic"`
    KeyPoints  []string `json:"key_points"`
    Difficulty int      `json:"difficulty"`
}
```

这就是你最后想拿到的结构化结果。

### 第二步：创建 JSON parser

```go
parser := schema.NewMessageJSONParser[LessonSummary](&schema.MessageJSONParseConfig{
    ParseFrom: schema.MessageParseFromContent,
})
```

这里最重要的是 `ParseFrom`。

它告诉 parser 去哪里找 JSON：

- `MessageParseFromContent`：从 `message.Content` 里取
- `MessageParseFromToolCall`：从工具调用参数里取

这一课先用最简单的 content 解析。

### 第三步：把 parser 放进 compose

```go
chain := compose.NewChain[*schema.Message, LessonSummary]()
chain.AppendLambda(compose.MessageParser(parser))
```

这一步的意义是：

- `MessageParser` 既可以单独用
- 也可以作为 `Lambda` 节点接进编排链路

### 第四步：执行解析

```go
summary, err := runner.Invoke(ctx, message)
```

输入一条 message，输出一个结构体。

这就是结构化解析最核心的价值。

## 5. 本课真正要记住的事

只记这四点：

1. `MessageJSONParser` 用来把 message 里的 JSON 解析成结构体
2. `ParseFrom` 决定从 content 还是 tool call 参数里取数据
3. 目标类型由泛型决定
4. `compose.MessageParser(parser)` 可以把它接进 Chain

## 6. 下一课讲什么

下一课讲 `ChatModelAgent`。

因为到这里你已经具备了：

- 模型调用
- 工具调用
- 编排
- 结构化解析

下一步就是把这些能力放进更高层的 agent 运行框架里。

## 7. 官方资料

- ChatModel Guide: https://www.cloudwego.io/docs/eino/core_modules/components/chat_model_guide/
- Eino 仓库: https://github.com/cloudwego/eino
