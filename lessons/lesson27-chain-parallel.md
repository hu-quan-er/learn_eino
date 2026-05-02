# 第 27 课：学会 Chain Parallel

这一课只讲一件事：在 `Chain` 里做并行扇出。

学完这一课，你要记住：

```text
Chain -> AppendParallel -> map output
```

## 1. 这一课解决什么问题

第 24 课你已经学过 graph 的并行汇聚。

但有些时候，你并不想单独建一个 graph，只是想在 chain 的某一个位置：

- 把同一份输入同时交给多个节点
- 再拿到一个汇总后的 map 继续往下走

这就是 `AppendParallel(...)` 的用法。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson27-chain-parallel/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson27-chain-parallel
```

demo 会打印：

1. 并行聚合后的最终字符串
2. 总耗时 `elapsed`

## 4. 这节课最关键的代码

### 第一步：先准备一个 Parallel 容器

```go
parallel := compose.NewParallel()
parallel.
    AddLambda("outline", ...).
    AddLambda("keywords", ...)
```

这里的 `"outline"`、`"keywords"`` 就是输出 key。

### 第二步：把 Parallel 插进 chain

```go
chain.
    AppendLambda(...).
    AppendParallel(parallel).
    AppendLambda(...)
```

这表示：

- 上一个节点的输出
- 会同时喂给 parallel 里的每个子节点

### 第三步：下游拿到的是一个 map

```go
AppendLambda(func(ctx context.Context, input map[string]any) (string, error) {
    outline, _ := input["outline"].(string)
    keywords, _ := input["keywords"].(string)
    ...
})
```

这是本课核心。

`AppendParallel(...)` 之后，下游一般接到的是 `map[string]any`。

## 5. 本课真正要记住的事

1. `NewParallel()` 是 chain 并行扇出的入口
2. 每个子节点都要有自己的 output key
3. `AppendParallel(...)` 后通常会得到一个 `map[string]any`
4. 它适合“线性流程中间插一个并行段”
5. 如果整体已经是复杂 DAG，再优先考虑直接写 graph

## 6. 它和第 24 课 Graph 并行汇聚的区别

`Chain Parallel` 更像：

- 在线性链路里插一个并行步骤

`Graph 并行汇聚` 更像：

- 整体流程本来就是 DAG

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
