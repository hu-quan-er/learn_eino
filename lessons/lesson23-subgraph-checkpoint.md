# 第 23 课：给 SubGraph 加上 Checkpoint

这一课只讲一件事：理解“外层 graph 正常推进，内层 sub-graph 暂停并恢复”的最小闭环。

学完这一课，你要记住：

```text
outer Graph -> subGraph -> interrupt -> ResumeWithData
```

## 1. 这一课解决什么问题

第 15 课你已经学过：

- 整个 workflow 可以 checkpoint

但真实应用更常见的情况是：

- 外层流程继续存在
- 真正需要暂停的是里面某个子流程
- 恢复时还要把额外状态塞回去

这节课讲的就是这个层级。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson23-subgraph-checkpoint/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson23-subgraph-checkpoint
```

你会看到两段输出：

1. 第一次运行时，在 sub-graph 里中断
2. 恢复后，把 `Reviewer=Alice` 填回去，继续跑完

## 4. 这节课最关键的代码

### 第一步：先给 sub-graph 配 local state

```go
subGraph := compose.NewGraph[string, string](
    compose.WithGenLocalState(func(ctx context.Context) *ReviewState {
        return &ReviewState{}
    }),
)
```

这里的 `ReviewState` 会成为 sub-graph 的局部状态。

如果这个 state 会进入 checkpoint，还要记得提前注册：

```go
schema.Register[ReviewState]()
```

### 第二步：把中断点放到 sub-graph 内部

```go
graph.AddGraphNode("content_pipeline", subGraph, compose.WithGraphCompileOptions(
    compose.WithGraphName("lesson23_content_pipeline"),
    compose.WithInterruptAfterNodes([]string{"draft"}),
))
```

注意这个位置。

中断配置不是放在外层节点名上，而是通过 `WithGraphCompileOptions(...)` 传给内层 sub-graph。

### 第三步：恢复时把状态塞回 root cause

```go
resumeCtx := compose.ResumeWithData(ctx, rootInterruptID(info), &ReviewState{Reviewer: "Alice"})
output, err := runner.Invoke(resumeCtx, "第 23 课：嵌套 Checkpoint", compose.WithCheckPointID(checkpointID))
```

这是本课最关键的一步。

恢复时不是重新建一个 graph，也不是手改内部对象，而是：

1. 从 interrupt info 里找到 root cause
2. 用 `ResumeWithData(...)` 把恢复数据挂到 context
3. 再用相同 checkpoint ID 继续执行

## 5. 本课真正要记住的事

1. checkpoint 不只可以用在顶层流程
2. sub-graph 也能单独中断和恢复
3. `WithGraphCompileOptions(...)` 是给嵌套 graph 传编译参数的关键入口
4. `ResumeWithData(...)` 适合把审批、人审结果、外部状态回填到断点位置
5. 这是做“人机混合流程”时非常常见的一种写法

## 6. 什么时候会用到它

- 文稿写到一半，需要人工审批
- 工具子流程需要等待外部系统结果
- 某个局部步骤需要人类输入再继续

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
