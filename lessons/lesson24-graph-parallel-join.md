# 第 24 课：学会 Graph 并行汇聚

这一课只讲一件事：理解 Graph 里的“扇出后再汇聚”。

学完这一课，你要记住：

```text
START -> parallel nodes -> merge node -> END
```

## 1. 这一课解决什么问题

前面的 graph 例子很多是单线串行。

但一旦你要做下面这类事：

- 同时生成大纲
- 同时提取关键词
- 两边都完成后再合并

就会进入典型的 DAG 汇聚场景。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson24-graph-parallel-join/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson24-graph-parallel-join
```

除了最终结果，demo 还会打印一个 `elapsed`，用来让你直观看到它不是简单串行。

## 4. 这节课最关键的代码

### 第一步：两个并行节点都从 START 出发

```go
graph.AddEdge(compose.START, "outline")
graph.AddEdge(compose.START, "keywords")
```

这就形成了最基础的扇出。

### 第二步：两个节点都指向 merge

```go
graph.AddEdge("outline", "merge")
graph.AddEdge("keywords", "merge")
```

这表示：

- `merge` 有两个前驱
- 它需要处理来自两个分支的结果

### 第三步：编译时启用 `AllPredecessor`

```go
runner, err := graph.Compile(ctx,
    compose.WithNodeTriggerMode(compose.AllPredecessor),
)
```

这是本课核心。

`AllPredecessor` 的含义是：

- `merge` 只有在所有前驱都完成后才会触发

如果你不理解这一点，就很容易把 Graph 用错。

### 第四步：merge 节点的输入不再是单值

```go
func(ctx context.Context, input map[string]any) (string, error)
```

因为前面两个节点都用了 `WithOutputKey(...)`，所以 merge 收到的是一个 map。

## 5. 本课真正要记住的事

1. Graph 很适合表达扇出和汇聚
2. 汇聚节点通常要配合 `AllPredecessor`
3. 上游并行节点最好带 `WithOutputKey(...)`
4. merge 节点一般会收到 `map[string]any`
5. 这类写法比把所有逻辑堆进一个 lambda 更清晰

## 6. 什么时候会用到它

- 多路召回后统一重排
- 同时生成多个草稿片段再合并
- 并行做抽取、分类、总结，然后拼装最终结果

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
