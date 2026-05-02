# 第 39 课：学会 WithSkipTransferMessages

这一课只讲一件事：控制转交过程中生成的 transfer messages 要不要进入下游输入历史。

学完这一课，你要记住：

```text
Runner.Query(..., WithSkipTransferMessages()) -> child history without transfer helper messages
```

## 1. 这一课解决什么问题

默认情况下，agent 转交时会带上一些辅助上下文。

这通常是好事，但有时你想：

- 下游只看业务内容
- 不看 transfer 过程噪音

这时就可以用 `WithSkipTransferMessages()`。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson39-skip-transfer-messages/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson39-skip-transfer-messages
```

你会看到两组结果：

1. 默认历史
2. 跳过 transfer messages 后的历史

## 4. 这节课最关键的代码

### 第一种：默认执行

```go
runner.Query(ctx, "lesson39 要强调什么")
```

### 第二种：带 `WithSkipTransferMessages`

```go
runner.Query(ctx, "lesson39 要强调什么", adk.WithSkipTransferMessages())
```

两次运行的 agent 代码完全一样，差别只在 runtime option。

## 5. 本课真正要记住的事

1. `WithSkipTransferMessages()` 是运行时选项
2. 它不会改 agent 实现，只会改一次 run 的 history 传递方式
3. 当 transfer 上下文太吵时，这个选项很有用
4. 它和 `HistoryRewriter` 不冲突，两个可以一起用
5. 生产里要根据下游 agent 的输入敏感度决定是否开启

## 6. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
