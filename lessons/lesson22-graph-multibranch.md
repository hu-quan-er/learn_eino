# 第 22 课：学会 Graph MultiBranch

这一课只讲一件事：让一个 Graph 在同一次运行里，同时选中多个分支。

学完这一课，你要记住：

```text
Graph -> NewGraphMultiBranch -> selected nodes -> map output
```

## 1. 这一课解决什么问题

前面的 branch 例子大多是：

- 二选一
- 选一个节点继续走

但实际业务里经常不是“选一个”，而是：

- 同时执行多个分支
- 最后把多个结果一起收回来

这就是 `NewGraphMultiBranch(...)` 的作用。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson22-graph-multibranch/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson22-graph-multibranch
```

demo 会跑 3 个 case：

1. `Eino 实战`
2. `Eino 答疑`
3. `Eino 入门`

## 4. 这节课最关键的代码

### 第一步：给每个分支节点加上 output key

```go
graph.AddLambdaNode("summary", ..., compose.WithOutputKey("summary"))
graph.AddLambdaNode("examples", ..., compose.WithOutputKey("examples"))
graph.AddLambdaNode("faq", ..., compose.WithOutputKey("faq"))
```

因为最后会有多个节点同时产出结果，所以不能再只靠一个单值输出。这里要让每个节点把结果写进自己的 key。

### 第二步：返回一个 `map[string]bool`

```go
compose.NewGraphMultiBranch(func(ctx context.Context, input string) (map[string]bool, error) {
    selected := map[string]bool{"summary": true}
    if strings.Contains(input, "实战") {
        selected["examples"] = true
    }
    if strings.Contains(input, "答疑") {
        selected["faq"] = true
    }
    return selected, nil
}, ...)
```

这一步是本课核心。

普通 branch 返回一个节点名；
`MultiBranch` 返回的是“一组节点名”。

### 第三步：最终结果会是 map

```go
output, err := runner.Invoke(ctx, "Eino 实战")
```

返回值不是单个字符串，而是：

```go
map[string]any{
    "summary":  "...",
    "examples": "...",
}
```

只有被选中的节点才会出现在结果里。

## 5. 本课真正要记住的事

1. `NewGraphMultiBranch` 可以一次选中多个节点
2. 多分支结果通常要配合 `WithOutputKey(...)`
3. 输出只包含本次实际走到的分支
4. 这很适合“按条件拼装多个处理块”
5. 它和 workflow 的单路由思路不一样，Graph 更适合这种 DAG 式扇出

## 6. 什么时候会用到它

- 一个问题同时需要“摘要 + 示例”
- 一个请求同时触发“主处理 + FAQ 生成”
- 根据标签决定要跑哪些后处理模块

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
