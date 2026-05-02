# 第 11 课：让 Tool 也能流式输出

这一课只讲一件事：把前面学过的 `Tool` 升级成可以一段一段返回结果的 `StreamableTool`。

学完这一课，你要记住：

```text
StreamableTool -> ToolsNode.Stream -> ToolMessage chunks
```

## 1. 这一课解决什么问题

前面的 tool demo 都是“一次性返回完整结果”。

但现实里很多工具结果本身就是分段产生的，比如：

- 搜索过程中的中间进度
- 长文本生成
- 长日志输出
- 实时执行状态

这时你就不该继续用一次性 `InvokableRun`，而应该改成流式输出。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson11-streamable-tool/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson11-streamable-tool
```

你会看到工具结果被分成多段打印，最后再拼成一条完整 `tool message`。

## 4. 这节课最关键的代码

### 第一步：用 `InferStreamTool` 定义流式工具

```go
searchTool, err := toolutils.InferStreamTool(
    "search_docs",
    "stream mock search result chunks",
    streamSearchDocs,
)
```

这里和第 04 课最大的区别是：

- 不再用 `InferTool`
- 改成用 `InferStreamTool`

### 第二步：工具函数返回 `*schema.StreamReader[string]`

```go
func streamSearchDocs(ctx context.Context, input *SearchInput) (*schema.StreamReader[string], error)
```

这个签名说明：

- 输入还是正常结构体
- 输出变成 stream reader

demo 里用 `schema.Pipe(...)` 手动构造了一个最小流。

### 第三步：用 `ToolsNode.Stream(...)` 执行

```go
stream, err := toolsNode.Stream(ctx, assistantMessage)
```

注意这里也不是 `Invoke(...)` 了，而是 `Stream(...)`。

### 第四步：逐 chunk 读取

```go
messages, err := stream.Recv()
```

每次读到的是一批 `[]*schema.Message`，因为 `ToolsNode` 允许一次处理多个 tool call。

本课 demo 只有一个 tool call，所以你看到的每个 chunk 实际上都只有一条 `tool message`。

### 第五步：最后把 chunk 合并回来

```go
merged, err := schema.ConcatMessageArray(chunks)
```

这一步很重要。

因为上游经常同时需要两种结果：

- 实时 chunk
- 完整消息

## 5. 本课真正要记住的事

1. 流式工具定义入口是 `InferStreamTool`
2. 工具函数返回值要变成 `*schema.StreamReader[T]`
3. 执行入口从 `ToolsNode.Invoke(...)` 变成 `ToolsNode.Stream(...)`
4. 读取时仍然要处理 `io.EOF`
5. 最后可以用 `schema.ConcatMessageArray(...)` 把多个 chunk 合并回来

## 6. 和第 02 课有什么关系

第 02 课学的是：

- 模型怎么流式输出

第 11 课学的是：

- 工具怎么流式输出

这两个概念很像，但位置不同：

- 第 02 课是在 `ChatModel`
- 第 11 课是在 `Tool`

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
