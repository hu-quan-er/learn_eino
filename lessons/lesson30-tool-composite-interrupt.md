# 第 30 课：学会 Tool CompositeInterrupt

这一课只讲一件事：当一个 tool 内部又调用了 graph 时，怎么把内部中断正确向外传。

学完这一课，你要记住：

```text
tool -> inner graph -> compose.Interrupt -> tool.CompositeInterrupt
```

## 1. 这一课解决什么问题

前一课你学的是：

- tool 自己中断自己

但还有一种更常见的进阶场景：

- tool 内部其实还包了一层 graph
- 真正中断的是 graph 里面的节点
- 外层 tool 需要把这个中断“包装”后继续往上抛

这就是 `tool.CompositeInterrupt(...)` 的职责。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson30-tool-composite-interrupt/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson30-tool-composite-interrupt
```

你会看到：

1. 第一次运行时，根因中断来自 inner graph
2. parent info 里能看到 wrapper tool 自己的中断信息
3. 恢复后，inner graph 正常继续

## 4. 这节课最关键的代码

### 第一步：inner graph 自己发中断

```go
if !wasInterrupted {
    return "", compose.Interrupt(ctx, "inner graph needs resume")
}
```

这里中断发生在 graph 内部节点，不在 tool 本体。

### 第二步：wrapper tool 捕获并包装

```go
result, err := t.compiledGraph.Invoke(ctx, argumentsInJSON)
if err != nil {
    if _, ok := compose.ExtractInterruptInfo(err); ok {
        return "", tool.CompositeInterrupt(ctx, "wrapper tool interrupt", nil, err)
    }
    return "", err
}
```

这是本课核心。

`CompositeInterrupt(...)` 做的事情是：

- 保留内部真正的 root cause
- 同时把当前 tool 自己也挂到 parent 链里

### 第三步：恢复时直接恢复 root cause

```go
resumeCtx := compose.Resume(ctx, root.ID)
output, err := runner.Invoke(resumeCtx, input, compose.WithCheckPointID(checkpointID))
```

恢复目标是最底层真正中断的那个点，而不是随便恢复外层某个节点。

## 5. 本课真正要记住的事

1. `CompositeInterrupt` 适合“组件里又包了组件”的情况
2. root cause 仍然是最里面真正中断的那个点
3. 外层 wrapper 的中断信息会出现在 parent 链里
4. 这能保证恢复路径既准确，又保留外层语义
5. 当 tool 内部越来越像一个子系统时，这种写法很重要

## 6. 什么时候会用到它

- 一个 tool 内部跑子 graph
- 一个 tool 内部再调其他可中断组件
- 想把复杂内部流程封成单个 tool 暴露给上层

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
