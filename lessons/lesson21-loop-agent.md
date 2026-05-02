# 第 21 课：学会 LoopAgent

这一课只讲一件事：理解 Eino ADK 里“重复执行同一组 agent”的最小写法。

学完这一课，你要记住：

```text
LoopAgent -> iterations -> BreakLoopAction
```

## 1. 这一课解决什么问题

前面的 `SequentialAgent` 和 `ParallelAgent` 解决的是：

- 顺序执行
- 并行执行

但还有一种很常见的模式：

- 重复执行，直到次数耗尽
- 或者在中途主动跳出

这就是 `LoopAgent` 的职责。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson21-loop-agent/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson21-loop-agent
```

你会看到两组输出：

1. `plain loop`
2. `break loop`

## 4. 这节课最关键的代码

### 第一步：创建一个普通 LoopAgent

```go
plainLoop, err := adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
    Name:        "plain_loop",
    Description: "run the same agent for a fixed number of iterations",
    SubAgents:   []adk.Agent{...},
    MaxIterations: 3,
})
```

这表示：

- 子 agent 每轮都执行一次
- 最多执行 3 轮

### 第二步：在子 agent 内部发出 BreakLoopAction

```go
if a.breakAfter > 0 && a.current >= a.breakAfter {
    event.Action = adk.NewBreakLoopAction(a.name)
}
```

这一步是本课核心。

`BreakLoopAction` 不是给模型看的，而是给框架看的控制信号。它的作用是：

- 当前轮结束后
- 提前停止后续迭代

### 第三步：观察框架补上的状态

```go
event.Action.BreakLoop.Done
event.Action.BreakLoop.CurrentIterations
```

`LoopAgent` 在真正处理完 break 之后，会把这些字段补上。也就是说，最终你看到的事件里会带上：

- `Done=true`
- `CurrentIterations=<在哪一轮退出>`

## 5. 本课真正要记住的事

1. `LoopAgent` 适合“重复尝试”或“多轮处理”的场景
2. `MaxIterations` 是最基础的止损手段
3. `BreakLoopAction` 允许子 agent 主动终止循环
4. `BreakLoopAction` 是程序控制信号，不是普通业务输出
5. 读事件时，仍然还是 `Runner -> AgentEvent`

## 6. 什么时候会用到 LoopAgent

- 计划执行器一轮轮拆任务
- 审批流反复修订
- 多轮自检，直到满足某个条件

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
