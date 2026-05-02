# 第 36 课：用 SetSubAgents 跑通自定义 Agent 转交

这一课只讲一件事：不用 `ChatModelAgent`，直接让一个自定义 agent 转交给另一个自定义 agent。

学完这一课，你要记住：

```text
SetSubAgents -> NewTransferToAgentAction -> child agent receives rewritten history
```

## 1. 这一课解决什么问题

前面你已经会：

- 单独跑一个 agent
- 跑 workflow agent

这一课开始看更底层的一层：

- 多个普通 `Agent` 之间怎么直接转交

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson36-set-subagents-transfer/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson36-set-subagents-transfer
```

## 4. 这节课最关键的代码

### 第一步：先把子 agent 挂上去

```go
root, err := adk.SetSubAgents(ctx, router, []adk.Agent{writer})
```

这一步之后，`router` 才真的知道：

- 哪些名字是自己可以转交的目标

### 第二步：父 agent 发出转交动作

```go
Action: adk.NewTransferToAgentAction("writer_agent")
```

这不是普通日志，也不是普通 message。

这是 ADK 控制流动作。

### 第三步：子 agent 收到的不是“原始裸输入”

demo 里 `writer_agent` 会把自己实际收到的 `input.Messages` 打印出来。

你会看到：

- 原始用户问题
- 父 agent 的消息重写结果
- transfer 相关上下文

## 5. 本课真正要记住的事

1. `SetSubAgents(...)` 是自定义 agent 转交的基础
2. 转交动作靠 `NewTransferToAgentAction(...)`
3. 子 agent 收到的是“重写后的 history”，不是只收到最后一句话
4. ADK 的多 agent 不一定要从 `ChatModelAgent` 开始
5. 先把这条底层链路看懂，再看更高层封装会更清楚

## 6. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
