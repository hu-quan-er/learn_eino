# 第 37 课：学会 ChatModelAgent 的 OutputKey

这一课只讲一件事：让 `ChatModelAgent` 自动把输出写进 ADK session。

学完这一课，你要记住：

```text
ChatModelAgent(OutputKey) -> AddSessionValue(auto) -> next agent reads session
```

## 1. 这一课解决什么问题

有时候你希望：

- 一个 agent 先产出中间结果
- 后面的 agent 直接复用

如果这个中间结果不适合再从 message history 里硬解析，`OutputKey` 很合适。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson37-chatmodel-output-key/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson37-chatmodel-output-key
```

## 4. 这节课最关键的代码

### 第一步：给 ChatModelAgent 配 `OutputKey`

```go
OutputKey: "draft_output",
```

含义是：

- agent 最终输出的 `message.Content`
- 会自动写到 session 的 `draft_output`

### 第二步：后一个 agent 直接从 session 读

```go
value, _ := adk.GetSessionValue(ctx, "draft_output")
```

这里不需要自己额外 `AddSessionValue(...)`，因为框架已经替你写进去了。

### 第三步：把它放进顺序 agent 里最容易观察

```go
SubAgents: []adk.Agent{chatAgent, reader},
```

这样你能直接看到：

- 前一个 agent 产出 draft
- 后一个 agent 读到 draft

## 5. 本课真正要记住的事

1. `OutputKey` 是 `ChatModelAgent` 自带的 session 输出能力
2. 它适合存“后续逻辑还要用”的模型结果
3. 后续 agent 用 `GetSessionValue(...)` 直接读
4. 这比重新从 message history 里抽值更稳定
5. `OutputKey` 和第 34 课的手动 session values 可以组合使用

## 6. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
