# 第 02 课：学会流式输出 Stream

这一课只讲一件事：把上一课的一次性 `Generate`，换成流式 `Stream`。

学完这一课，你要记住的最小调用路径是：

```text
创建 ChatModel -> 调用 Stream -> 循环 Recv -> 遇到 io.EOF 结束 -> Close
```

## 1. 这一课解决什么问题

上一课的 `Generate` 很直接：

- 发请求
- 等模型完整生成
- 一次性拿到最终答案

这适合很多场景，但有一个明显问题：

- 用户必须等完整答案生成完，界面才有内容

所以第二课要学 `Stream`。它适合：

- 聊天窗口逐字输出
- 长回答降低等待感
- 需要边生成边展示的交互

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson02-stream/main.go`

这个 demo 和第一课的模型初始化方式完全一致，只把“取结果”的方式换成了流式读取。

## 3. 如何运行

先进入目录并准备环境变量：

```bash
cd /Users/seaya/huquan.95/project-4/eino
export OPENAI_API_KEY="你的密钥"
export OPENAI_MODEL="gpt-4o-mini"
export OPENAI_BASE_URL="你的兼容接口地址"
```

执行：

```bash
go run ./cmd/lesson02-stream
```

如果配置正常，你会看到两段输出：

```text
assistant (streaming):
1. ...
2. ...
3. ...

assistant (merged):
1. ...
2. ...
3. ...

received 12 chunks
```

说明：

- 第一段是边收边打印的流式输出
- 第二段是把所有 chunk 合并后得到的完整答案
- `received 12 chunks` 只是示意，实际数量取决于模型服务

## 4. 先看最关键的变化

第一课是：

```go
resp, err := chatModel.Generate(ctx, messages, model.WithTemperature(0.2))
```

第二课变成：

```go
stream, err := chatModel.Stream(ctx, messages, model.WithTemperature(0.2))
```

核心差别只有一句话：

- `Generate` 返回完整消息
- `Stream` 返回消息流读取器

所以 `Stream` 之后，你不能直接 `resp.Content`，而是要不断调用 `Recv()`。

## 5. 代码怎么读

### 第一步：拿到流

```go
stream, err := chatModel.Stream(ctx, messages, model.WithTemperature(0.2))
if err != nil {
    log.Fatalf("stream failed: %v", err)
}
defer stream.Close()
```

这里要注意两件事：

1. `Stream(...)` 返回的是 `*schema.StreamReader[*schema.Message]`
2. 读完后要 `Close()`

这一点不是可有可无。Eino 当前版本源码的注释明确建议：用了 `Recv()` 之后，记得 `Close()`。

### 第二步：循环 Recv

```go
for {
    chunk, err := stream.Recv()
    if errors.Is(err, io.EOF) {
        break
    }
    if err != nil {
        log.Fatalf("stream recv failed: %v", err)
    }

    chunks = append(chunks, chunk)
    fmt.Print(chunk.Content)
}
```

这是第二课最核心的代码块。

你可以把它理解成：

- `Recv()` 取下一个 chunk
- 如果遇到 `io.EOF`，说明流正常结束
- 如果遇到别的错误，说明这次流式调用异常
- 每个 `chunk` 本质上仍然是 `*schema.Message`

这里最容易犯的错有两个：

1. 把 `io.EOF` 当成错误处理掉
2. 忘记 `Close()`

在流式编程里，`io.EOF` 在这里表示“正常读完了”，不是失败。

### 第三步：边收边打印

```go
fmt.Print(chunk.Content)
```

这就是流式输出看起来“一个字一个字往外冒”的原因。

因为每收到一个 chunk，就立刻输出一部分内容，而不是等全部完成。

### 第四步：把所有 chunk 拼回完整消息

```go
fullMessage, err := schema.ConcatMessages(chunks)
```

为什么这里不用 `schema.ConcatMessageStream(stream)`？

因为我们已经在上面的 `Recv()` 循环里把流读掉了。

所以这一课故意这样设计：

- 一边 `Recv()` 一边打印
- 同时把 chunk 存到切片里
- 最后用 `schema.ConcatMessages(chunks)` 拼回完整答案

这样你能同时理解两件事：

1. 流式输出是怎么实时展示的
2. stream 的每个 chunk 本质上仍然是 message

## 6. Generate 和 Stream 到底怎么选

你现在先按这个简单标准记：

用 `Generate`：

- 后台任务
- 不需要逐步展示
- 只关心最终答案

用 `Stream`：

- 对话 UI
- 回答比较长
- 想尽快让用户看到内容

## 7. 本课真正要记住的事

只记这五点：

1. `Stream` 返回的不是完整结果，而是 `StreamReader`
2. 要循环调用 `Recv()` 才能拿到内容
3. `io.EOF` 在这里表示正常结束
4. 读完后要 `Close()`
5. 流里的每个 chunk 仍然是 `schema.Message`

## 8. 你现在可以自己做的练习

建议你立刻改三处再跑一遍：

1. 把用户问题改成长一点，观察 chunk 数量变化
2. 把 `fmt.Print(chunk.Content)` 改成 `fmt.Printf("[%q]", chunk.Content)`，观察 chunk 边界
3. 删掉 `chunks = append(...)` 和合并逻辑，只保留纯流式打印

## 9. 下一课讲什么

第二课之后，最自然的下一步不是 Agent，而是 `Prompt`。

原因很直接：

- 你已经会了 `Generate`
- 也会了 `Stream`
- 下一步应该学“怎么更稳定地组织输入”

所以第三课我建议讲 `Prompt Template`，把写死的消息改成可复用模板。

## 10. 官方资料

本课内容参考的是 Eino 官方一手资料和当前版本源码：

- Quick Start: https://www.cloudwego.io/docs/eino/quick_start/
- ChatModel Guide: https://www.cloudwego.io/docs/eino/core_modules/components/chat_model_guide/
- Stream Programming Essentials: https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/stream_programming_essentials/
- Eino 仓库: https://github.com/cloudwego/eino
