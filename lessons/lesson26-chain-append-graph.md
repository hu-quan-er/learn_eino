# 第 26 课：学会 Chain 里挂一个 Graph

这一课只讲一件事：理解 `Chain` 不是只能串 lambda，也可以直接把一个完整 `Graph` 当成中间节点。

学完这一课，你要记住：

```text
Chain -> AppendLambda -> AppendGraph -> AppendLambda
```

## 1. 这一课解决什么问题

前面你已经学过：

- `Graph` 适合表达局部流程
- `Workflow` 适合表达外层总控

但在很多场景里，还有一种很实用的写法：

- 外层只想要线性调用体验
- 中间某一步又不想只写成一个 lambda

这时候就可以在 `Chain` 里直接挂一个 `Graph`。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson26-chain-append-graph/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson26-chain-append-graph
```

## 4. 这节课最关键的代码

### 第一步：先写一个子 graph

```go
func buildDraftGraph() *compose.Graph[string, string] { ... }
```

这里的子 graph 只负责：

- `outline`
- `write`

也就是“提纲 -> 初稿”。

### 第二步：在 chain 里挂进去

```go
chain.
    AppendLambda(...).
    AppendGraph(buildDraftGraph()).
    AppendLambda(...)
```

这里是本课核心。

`AppendGraph(...)` 的意思不是“切换到 graph 模式”，而是：

- 把这个 graph 当成 chain 里的一个普通节点
- 输入接上游
- 输出接下游

### 第三步：顶层执行入口还是 chain

```go
runner, err := chain.Compile(ctx)
output, err := runner.Invoke(ctx, "...")
```

外部调用时，并不会因为中间嵌了 graph 就改成另一套运行方式。

## 5. 本课真正要记住的事

1. `Chain` 可以嵌套 `Graph`
2. `AppendGraph(...)` 适合把一段局部流程封成一个链路节点
3. 外部执行入口仍然是 `chain.Compile().Invoke(...)`
4. 当你只想保留“线性阅读感”时，这种写法很合适
5. 这是把复杂度藏在内部的一种常见手段

## 6. 什么时候会用到它

- 线性主流程里嵌一个小型草稿流水线
- 线性主流程里嵌一个检索子图
- 想保留 chain 风格，但中间有 2 到 3 步局部处理

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
