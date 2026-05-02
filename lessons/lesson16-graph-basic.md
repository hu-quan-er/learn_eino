# 第 16 课：学会 Graph 基础编排

这一课只讲一件事：第一次直接使用 `Graph`，理解最底层的节点和边是怎么连起来的。

学完这一课，你要记住：

```text
node -> AddEdge -> Compile -> Invoke
```

## 1. 这一课解决什么问题

前面你已经学过：

- `Chain`
- `Workflow`

但 Eino 里更底层的编排单元其实是 `Graph`。

`Graph` 的特点是：

- 你自己显式加节点
- 你自己显式连边
- 抽象更低，也更灵活

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson16-graph-basic/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson16-graph-basic
```

## 4. 这节课最关键的代码

### 第一步：创建 Graph

```go
graph := compose.NewGraph[string, string]()
```

这里的泛型表示：

- 图的输入是 `string`
- 图的输出也是 `string`

### 第二步：加节点

```go
graph.AddLambdaNode("normalize", ...)
graph.AddLambdaNode("reply", ...)
```

和 workflow 不同，这里你直接面对图里的节点。

### 第三步：显式连边

```go
graph.AddEdge(compose.START, "normalize")
graph.AddEdge("normalize", "reply")
graph.AddEdge("reply", compose.END)
```

这是这一课最核心的地方。

`Graph` 的执行顺序不是靠你“声明输入依赖”推出来的，而是靠边明确指定。

### 第四步：编译和执行

```go
runner, err := graph.Compile(ctx, compose.WithGraphName("lesson16_graph"))
output, err := runner.Invoke(ctx, "  什么情况下要直接使用 Graph？  ")
```

## 5. 本课真正要记住的事

1. `Graph` 是更底层的编排能力
2. `AddLambdaNode(...)` 是加节点
3. `AddEdge(...)` 是连边
4. `START` 和 `END` 是保留节点
5. `Compile(...)` 之后才可以 `Invoke(...)`

## 6. 和 Workflow 的区别

你可以先这样理解：

- `Workflow` 更适合声明式数据编排
- `Graph` 更适合你想精确控制执行路径的时候

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
