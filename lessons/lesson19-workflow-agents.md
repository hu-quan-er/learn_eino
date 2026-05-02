# 第 19 课：学会 SequentialAgent 和 ParallelAgent

这一课只讲一件事：理解 Eino ADK 里多 agent 编排的两种基础模式。

学完这一课，你要记住：

```text
SequentialAgent / ParallelAgent -> Runner -> AgentEvent
```

## 1. 这一课解决什么问题

当一个任务开始拆分成多个角色时，你通常会遇到两种模式：

- 先后执行
- 并行执行

Eino 已经把这两种模式抽象成了：

- `SequentialAgent`
- `ParallelAgent`

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson19-workflow-agents/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson19-workflow-agents
```

你会看到两组输出：

1. sequential events
2. parallel events

## 4. 这节课最关键的代码

### 第一步：创建子 agent

```go
researchAgent := &TimedAgent{...}
writeAgent := &TimedAgent{...}
```

本课用了本地 mock agent，目的是把注意力集中在“编排模式”而不是模型调用上。

### 第二步：创建 SequentialAgent

```go
sequentialAgent, err := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
    Name:        "sequential_demo",
    Description: "run sub agents one by one",
    SubAgents:   []adk.Agent{researchAgent, writeAgent},
})
```

它的含义很直接：

- 子 agent 按顺序执行

### 第三步：创建 ParallelAgent

```go
parallelAgent, err := adk.NewParallelAgent(ctx, &adk.ParallelAgentConfig{
    Name:        "parallel_demo",
    Description: "run sub agents in parallel",
    SubAgents:   []adk.Agent{searchAgent, reviewAgent},
})
```

它表示：

- 子 agent 可以并发执行

### 第四步：统一还是通过 Runner 读取事件

```go
consumeEvents(adk.NewRunner(ctx, adk.RunnerConfig{Agent: sequentialAgent}).Query(ctx, "..."))
consumeEvents(adk.NewRunner(ctx, adk.RunnerConfig{Agent: parallelAgent}).Query(ctx, "..."))
```

这一点很重要。

虽然 agent 结构变复杂了，但顶层消费方式没有变，仍然是 `Runner -> AgentEvent`。

## 5. 本课真正要记住的事

1. `SequentialAgent` 适合流水线式多角色执行
2. `ParallelAgent` 适合可并发拆分的子任务
3. 两者都还是 `Agent`
4. 顶层执行入口仍然是 `Runner`
5. 多 agent 编排不一定先依赖真实模型，可以先用本地 mock agent 把流程走通

## 6. 怎么选择 Sequential 还是 Parallel

用 `SequentialAgent`：

- 后一个角色依赖前一个角色结果

用 `ParallelAgent`：

- 各个角色可以独立做事
- 最后再汇总

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
