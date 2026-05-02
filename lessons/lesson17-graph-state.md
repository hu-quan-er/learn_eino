# 第 17 课：在 Graph 里共享 State

这一课只讲一件事：让多个节点通过本地 state 共享中间信息。

学完这一课，你要记住：

```text
WithGenLocalState -> WithStatePostHandler -> WithStatePreHandler
```

## 1. 这一课解决什么问题

有些信息不适合直接塞到节点输入输出里来回传。

比如：

- 某一步算出来的中间结果
- 节点执行计数
- 临时上下文

这时候就可以用 graph local state。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson17-graph-state/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson17-graph-state
```

你会看到最终输出里包含：

- 上一步保存下来的 `normalized`
- state 里累计的 `steps`

## 4. 这节课最关键的代码

### 第一步：开启 local state

```go
graph := compose.NewGraph[string, string](
    compose.WithGenLocalState(func(ctx context.Context) *LessonState {
        return &LessonState{}
    }),
)
```

没有这一步，后面的 state handler 都不能用。

### 第二步：在 post handler 里写 state

```go
compose.WithStatePostHandler(func(ctx context.Context, output string, state *LessonState) (string, error) {
    state.Normalized = output
    state.StepCount++
    return output, nil
})
```

这表示：

- 节点正常产出 `output`
- 同时把 `output` 写进 state

### 第三步：在 pre handler 里读 state

```go
compose.WithStatePreHandler(func(ctx context.Context, input string, state *LessonState) (string, error) {
    state.StepCount++
    return fmt.Sprintf("[normalized=%s steps=%d] %s", state.Normalized, state.StepCount, input), nil
})
```

这一步说明：

- 下游节点执行前
- 可以先根据 state 改造输入

## 5. 本课真正要记住的事

1. `WithGenLocalState(...)` 负责声明图级 state
2. `WithStatePostHandler(...)` 常用来“写 state”
3. `WithStatePreHandler(...)` 常用来“读 state”
4. state 生命周期属于一次 graph run
5. state 适合共享中间信息，不适合替代正常输入输出设计

## 6. 什么时候该用 state

适合：

- 记录中间结果
- 记录步骤计数
- 暂存本次运行上下文

不适合：

- 所有业务字段都往 state 里塞
- 把 state 当成数据库

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
