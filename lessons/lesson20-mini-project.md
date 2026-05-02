# 第 20 课：做一个完整小项目骨架

这一课只讲一件事：把前面学过的编排能力拼成一个更像真实应用的最小骨架。

学完这一课，你要记住：

```text
Workflow -> nested Graph -> review -> summary
```

## 1. 这一课解决什么问题

学单个组件时，你看到的都是“局部能力”。

但项目真正落地时，通常会出现这种结构：

1. 先整理输入
2. 再进入某个子流程
3. 子流程产出结果
4. 外层流程继续审阅和汇总

这节课就是把它组合起来。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson20-mini-project/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson20-mini-project
```

## 4. 这节课最关键的代码

### 第一步：先写一个内层 Graph

```go
func buildDraftGraph() *compose.Graph[string, string] { ... }
```

这个内层 graph 只做一件事：

- `plan -> write`

也就是“提纲 -> 初稿”。

### 第二步：外层 Workflow 把子流程挂进来

```go
workflow.AddGraphNode(
    "draft_pipeline",
    buildDraftGraph(),
    compose.WithGraphCompileOptions(compose.WithGraphName("draft_pipeline")),
)
```

这里是本课最关键的点。

你不是只能在 workflow 里放 lambda 节点，也可以放一个完整 sub-graph。

### 第三步：前后再接业务节点

外层 workflow 里还有：

- `prepare_request`
- `review`
- `summary`

这就是比较真实的项目骨架思路：

- 外层做流程控制
- 内层做某个相对独立的子流水线

### 第四步：最终统一收口

```go
workflow.End().AddInput("review", compose.MapFields("draft", "draft"))
workflow.End().AddInput("summary", compose.MapFields("summary", "summary"))
```

## 5. 本课真正要记住的事

1. 一个 workflow 里可以嵌套 graph
2. `AddGraphNode(...)` 适合封装子流程
3. 外层流程负责总控，内层 graph 负责局部复杂度
4. 这是最接近真实项目结构的入门写法
5. 当 demo 变大时，下一步就是把这些 builder 函数拆到独立文件

## 6. 到第 20 课为止，你已经掌握了什么

你现在已经完成了从组件到应用骨架的一条主线：

1. ChatModel
2. Stream
3. Prompt
4. Tool
5. Chain
6. ToolsNode
7. Model + Tool 闭环
8. Workflow
9. MessageParser
10. ChatModelAgent
11. StreamableTool
12. Workflow Branch
13. Agent Stream
14. Agent Interrupt / Resume
15. Workflow Checkpoint
16. Graph
17. Graph State
18. AgentTool
19. Sequential / Parallel Agent
20. Mini Project Skeleton

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
