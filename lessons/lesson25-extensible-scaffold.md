# 第 25 课：做一个可扩展脚手架

这一课只讲一件事：把前面学过的 `Workflow`、`Graph`、`Agent` 拼成一个更接近项目结构的最小脚手架。

学完这一课，你要记住：

```text
Workflow(prep) -> Graph(draft) -> SequentialAgent(review) -> final package
```

## 1. 这一课解决什么问题

到前面为止，你已经学了很多单点能力。

真正做项目时，更合理的拆法通常是：

1. 用 workflow 处理外层业务步骤
2. 用 graph 封装局部流水线
3. 用 agent 承担角色化审阅
4. 最后统一拼装结果

这节课就是把这条线完整串起来。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson25-extensible-scaffold/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson25-extensible-scaffold
```

你会看到两部分：

1. `review events`
2. `extensible scaffold output`

## 4. 这节课最关键的代码

### 第一步：外层 workflow 只管总控

```go
workflow := compose.NewWorkflow[map[string]any, map[string]any]()
```

workflow 里只放：

- `prepare_request`
- `draft_pipeline`
- `collect_draft`

也就是：

- 整理输入
- 跑草稿子流程
- 收集结果

### 第二步：把 draft pipeline 封装成一个 graph

```go
workflow.AddGraphNode("draft_pipeline", buildDraftGraph(), ...)
```

这个内层 graph 负责：

- 生成提纲
- 生成案例
- 最后合并成初稿

所以它是一个“局部复杂度容器”，而不是总控器。

### 第三步：把 review team 独立成 agent

```go
reviewTeam, err := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{...})
```

这一步很重要。

不是所有东西都应该塞进 graph/workflow。

像“核对员 -> 编辑”这种角色链，更适合用 agent 表达。

### 第四步：主程序负责把三种能力粘起来

```go
workflowOutput, err := runner.Invoke(...)
reviewEvents, reviewSummary, err := runReviewTeam(ctx, reviewTeam, draft)
```

这就是本课最想让你看到的结构：

- 编排层面用 `compose`
- 角色协作用 `adk`
- 最终在应用入口处收口

## 5. 本课真正要记住的事

1. 外层 workflow 适合做业务总控
2. 内层 graph 适合封装一个局部流水线
3. agent 适合表达有角色含义的处理链
4. 项目入口处负责把这些能力组装起来
5. 这比“所有逻辑塞进一个大函数”更容易扩展

## 6. 到第 25 课为止，你已经掌握了什么

你现在已经完成了从组件到小型架构骨架的一条主线：

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
21. LoopAgent
22. Graph MultiBranch
23. SubGraph Checkpoint
24. Graph Parallel Join
25. Extensible Scaffold

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
