# 第 38 课：学会自定义 HistoryRewriter

这一课只讲一件事：控制一个 agent 在接手任务时到底看到什么历史。

学完这一课，你要记住：

```text
AgentWithOptions(..., WithHistoryRewriter(...)) -> custom input history
```

## 1. 这一课解决什么问题

默认 history rewrite 很通用，但不是所有系统都想保留完整上下文。

很多时候你更想要：

- 压缩历史
- 只保留关键上下文
- 把长链路改成更短的输入

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson38-history-rewriter/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson38-history-rewriter
```

你会看到两组结果：

1. default history
2. custom history rewriter

## 4. 这节课最关键的代码

### 第一步：把 child agent 用 `AgentWithOptions` 包一下

```go
rewrittenChild := adk.AgentWithOptions(ctx, child, adk.WithHistoryRewriter(compactHistory))
```

这表示：

- 不是改框架全局行为
- 只是给这个 agent 单独换一套 history 规则

### 第二步：自定义 `HistoryRewriter`

```go
func compactHistory(_ context.Context, entries []*adk.HistoryEntry) ([]adk.Message, error)
```

这里你拿到的是结构化历史条目，不是裸字符串。

你可以自己决定：

- 留哪些条目
- 丢哪些条目
- 如何压缩成新的 messages

## 5. 本课真正要记住的事

1. `WithHistoryRewriter(...)` 是控制子 agent 输入历史的关键接口
2. 它作用在 agent 级别，不是全局级别
3. 合理压缩 history，可以显著降低上下文噪音
4. 多 agent 系统里，history rewrite 往往和 transfer 一起设计
5. 这是生产里非常常见的定制点

## 6. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
