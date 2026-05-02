# Eino 学习目录

这个目录用来按章节学习 Eino。每一课都固定包含两部分：

1. 一个可以单独运行的 demo
2. 一份配套中文讲义

## 学习顺序

### 第 01 课：跑通第一个 ChatModel

- 代码：`cmd/lesson01-chatmodel/main.go`
- 讲义：`lessons/lesson01-chatmodel.md`
- 目标：理解 Eino 最基础的调用链
  `ChatModel -> Message -> Generate -> Response`

### 第 02 课：学会流式输出 Stream

- 代码：`cmd/lesson02-stream/main.go`
- 讲义：`lessons/lesson02-stream.md`
- 目标：理解 Eino 流式调用的最小闭环
  `ChatModel -> Stream -> Recv -> io.EOF -> Close`

### 第 03 课：学会 Prompt Template

- 代码：`cmd/lesson03-prompt-template/main.go`
- 讲义：`lessons/lesson03-prompt-template.md`
- 目标：理解 Prompt 的最小职责
  `variables -> ChatTemplate -> []*schema.Message`

### 第 04 课：学会定义一个 Tool

- 代码：`cmd/lesson04-tool/main.go`
- 讲义：`lessons/lesson04-tool.md`
- 目标：理解 Tool 的两个核心部分
  `ToolInfo(schema) + InvokableRun(execution)`

### 第 05 课：用 Chain 串起 Prompt 和 Model

- 代码：`cmd/lesson05-compose-chain/main.go`
- 讲义：`lessons/lesson05-compose-chain.md`
- 目标：理解 Eino 编排的最小闭环
  `ChatTemplate -> ChatModel -> Lambda -> output`

### 第 06 课：学会 ToolsNode 执行工具调用

- 代码：`cmd/lesson06-tools-node/main.go`
- 讲义：`lessons/lesson06-tools-node.md`
- 目标：理解工具执行节点的最小闭环
  `assistant tool call -> ToolsNode -> tool messages`

### 第 07 课：跑通模型调用工具的最小闭环

- 代码：`cmd/lesson07-model-with-tools/main.go`
- 讲义：`lessons/lesson07-model-with-tools.md`
- 目标：理解模型、工具和工具结果回填的完整顺序
  `model -> tool call -> ToolsNode -> model final answer`

### 第 08 课：学会 Workflow 基础编排

- 代码：`cmd/lesson08-workflow/main.go`
- 讲义：`lessons/lesson08-workflow.md`
- 目标：理解 Workflow 的依赖和数据映射
  `START -> nodes -> END`

### 第 09 课：把 Message 解析成结构体

- 代码：`cmd/lesson09-message-parser/main.go`
- 讲义：`lessons/lesson09-message-parser.md`
- 目标：理解结构化解析的最小路径
  `message JSON -> MessageParser -> Go struct`

### 第 10 课：跑通第一个 ChatModelAgent

- 代码：`cmd/lesson10-chatmodel-agent/main.go`
- 讲义：`lessons/lesson10-chatmodel-agent.md`
- 目标：理解 ADK agent 的最小运行方式
  `ChatModelAgent -> Runner -> AgentEvent`

### 第 11 课：让 Tool 也能流式输出

- 代码：`cmd/lesson11-streamable-tool/main.go`
- 讲义：`lessons/lesson11-streamable-tool.md`
- 目标：理解流式工具的执行路径
  `StreamableTool -> ToolsNode.Stream -> ToolMessage chunks`

### 第 12 课：学会 Workflow 分支路由

- 代码：`cmd/lesson12-workflow-branch/main.go`
- 讲义：`lessons/lesson12-workflow-branch.md`
- 目标：理解 workflow 的动态选路
  `START -> Branch -> selected node -> END`

### 第 13 课：跑通 Agent 的流式事件

- 代码：`cmd/lesson13-agent-stream/main.go`
- 讲义：`lessons/lesson13-agent-stream.md`
- 目标：理解 ADK 流式事件的最小读取方式
  `Runner(streaming) -> AgentEvent -> MessageStream`

### 第 14 课：学会 Agent 中断与恢复

- 代码：`cmd/lesson14-agent-interrupt-resume/main.go`
- 讲义：`lessons/lesson14-agent-interrupt-resume.md`
- 目标：理解 ADK 的中断恢复闭环
  `StatefulInterrupt -> checkpoint -> ResumeWithParams`

### 第 15 课：给 Workflow 加上 Checkpoint

- 代码：`cmd/lesson15-workflow-checkpoint/main.go`
- 讲义：`lessons/lesson15-workflow-checkpoint.md`
- 目标：理解流程级暂停与继续
  `Workflow -> CheckPointStore -> interrupt -> resume`

### 第 16 课：学会 Graph 基础编排

- 代码：`cmd/lesson16-graph-basic/main.go`
- 讲义：`lessons/lesson16-graph-basic.md`
- 目标：理解 Graph 的节点和边
  `node -> AddEdge -> Compile -> Invoke`

### 第 17 课：在 Graph 里共享 State

- 代码：`cmd/lesson17-graph-state/main.go`
- 讲义：`lessons/lesson17-graph-state.md`
- 目标：理解 graph local state 的最小用法
  `WithGenLocalState -> state handlers`

### 第 18 课：把 Agent 包装成 Tool

- 代码：`cmd/lesson18-agent-tool/main.go`
- 讲义：`lessons/lesson18-agent-tool.md`
- 目标：理解 agent 作为 tool 的复用方式
  `Agent -> NewAgentTool -> ToolsNode`

### 第 19 课：学会 SequentialAgent 和 ParallelAgent

- 代码：`cmd/lesson19-workflow-agents/main.go`
- 讲义：`lessons/lesson19-workflow-agents.md`
- 目标：理解多 agent 的顺序与并行编排
  `SequentialAgent / ParallelAgent -> Runner`

### 第 20 课：做一个完整小项目骨架

- 代码：`cmd/lesson20-mini-project/main.go`
- 讲义：`lessons/lesson20-mini-project.md`
- 目标：理解外层 workflow + 内层 graph 的组合方式
  `Workflow -> nested Graph -> summary`

### 第 21 课：学会 LoopAgent

- 代码：`cmd/lesson21-loop-agent/main.go`
- 讲义：`lessons/lesson21-loop-agent.md`
- 目标：理解循环 agent 的两种结束方式
  `LoopAgent -> MaxIterations / BreakLoopAction`

### 第 22 课：学会 Graph MultiBranch

- 代码：`cmd/lesson22-graph-multibranch/main.go`
- 讲义：`lessons/lesson22-graph-multibranch.md`
- 目标：理解 graph 的多分支扇出
  `Graph -> MultiBranch -> selected outputs`

### 第 23 课：给 SubGraph 加上 Checkpoint

- 代码：`cmd/lesson23-subgraph-checkpoint/main.go`
- 讲义：`lessons/lesson23-subgraph-checkpoint.md`
- 目标：理解嵌套 graph 的暂停与恢复
  `outer Graph -> subGraph -> interrupt -> resume`

### 第 24 课：学会 Graph 并行汇聚

- 代码：`cmd/lesson24-graph-parallel-join/main.go`
- 讲义：`lessons/lesson24-graph-parallel-join.md`
- 目标：理解 graph 的扇出和 fan-in
  `parallel nodes -> AllPredecessor -> merge`

### 第 25 课：做一个可扩展脚手架

- 代码：`cmd/lesson25-extensible-scaffold/main.go`
- 讲义：`lessons/lesson25-extensible-scaffold.md`
- 目标：理解 workflow、graph、agent 的组合分层
  `Workflow -> Graph -> Agent -> final package`

### 第 26 课：学会 Chain 里挂一个 Graph

- 代码：`cmd/lesson26-chain-append-graph/main.go`
- 讲义：`lessons/lesson26-chain-append-graph.md`
- 目标：理解 chain 对 graph 的嵌套方式
  `Chain -> AppendGraph -> downstream`

### 第 27 课：学会 Chain Parallel

- 代码：`cmd/lesson27-chain-parallel/main.go`
- 讲义：`lessons/lesson27-chain-parallel.md`
- 目标：理解 chain 中间的并行扇出
  `AppendParallel -> map output -> next node`

### 第 28 课：学会 Workflow 字段映射

- 代码：`cmd/lesson28-workflow-field-mapping/main.go`
- 讲义：`lessons/lesson28-workflow-field-mapping.md`
- 目标：理解 workflow 的三种常用映射
  `FromField / MapFields / ToField`

### 第 29 课：学会 Tool Interrupt / Resume

- 代码：`cmd/lesson29-tool-interrupt-resume/main.go`
- 讲义：`lessons/lesson29-tool-interrupt-resume.md`
- 目标：理解 tool 级暂停与继续
  `StatefulInterrupt -> GetInterruptState -> ResumeWithData`

### 第 30 课：学会 Tool CompositeInterrupt

- 代码：`cmd/lesson30-tool-composite-interrupt/main.go`
- 讲义：`lessons/lesson30-tool-composite-interrupt.md`
- 目标：理解 tool 包装子图时的中断透传
  `inner interrupt -> CompositeInterrupt -> resume root cause`

### 第 31 课：自己实现一个最小 Agent

- 代码：`cmd/lesson31-custom-agent/main.go`
- 讲义：`lessons/lesson31-custom-agent.md`
- 目标：理解自定义 agent 的最小接口面
  `Name / Description / Run -> Runner -> AgentEvent`

### 第 32 课：自己实现一个 ResumableAgent

- 代码：`cmd/lesson32-custom-resumable-agent/main.go`
- 讲义：`lessons/lesson32-custom-resumable-agent.md`
- 目标：理解自定义 agent 的中断恢复
  `Run -> StatefulInterrupt -> checkpoint -> Resume`

### 第 33 课：自己实现一个 Streaming Agent

- 代码：`cmd/lesson33-custom-streaming-agent/main.go`
- 讲义：`lessons/lesson33-custom-streaming-agent.md`
- 目标：理解自定义 agent 的流式消息输出
  `MessageVariant(streaming) -> MessageStream -> AgentEvent`

### 第 34 课：用 Session Values 共享 Agent 运行态

- 代码：`cmd/lesson34-agent-session-values/main.go`
- 讲义：`lessons/lesson34-agent-session-values.md`
- 目标：理解 ADK session values 的写入与读取
  `WithSessionValues -> AddSessionValue -> GetSessionValue`

### 第 35 课：学会 AgentTool 的高级输入模式

- 代码：`cmd/lesson35-agent-tool-advanced/main.go`
- 讲义：`lessons/lesson35-agent-tool-advanced.md`
- 目标：理解 agent tool 输入如何映射到子 agent
  `default request / custom schema / full chat history`

## 目录结构

```text
eino/
├── .env.example
├── go.mod
├── cmd/
│   ├── lesson01-chatmodel/
│   │   └── main.go
│   ├── lesson02-stream/
│   │   └── main.go
│   ├── lesson03-prompt-template/
│   │   └── main.go
│   ├── lesson04-tool/
│   │   └── main.go
│   ├── lesson05-compose-chain/
│   │   └── main.go
│   ├── lesson06-tools-node/
│   │   └── main.go
│   ├── lesson07-model-with-tools/
│   │   └── main.go
│   ├── lesson08-workflow/
│   │   └── main.go
│   ├── lesson09-message-parser/
│   │   └── main.go
│   ├── lesson10-chatmodel-agent/
│   │   └── main.go
│   ├── lesson11-streamable-tool/
│   │   └── main.go
│   ├── lesson12-workflow-branch/
│   │   └── main.go
│   ├── lesson13-agent-stream/
│   │   └── main.go
│   ├── lesson14-agent-interrupt-resume/
│   │   └── main.go
│   ├── lesson15-workflow-checkpoint/
│   │   └── main.go
│   ├── lesson16-graph-basic/
│   │   └── main.go
│   ├── lesson17-graph-state/
│   │   └── main.go
│   ├── lesson18-agent-tool/
│   │   └── main.go
│   ├── lesson19-workflow-agents/
│   │   └── main.go
│   ├── lesson20-mini-project/
│   │   └── main.go
│   ├── lesson21-loop-agent/
│   │   └── main.go
│   ├── lesson22-graph-multibranch/
│   │   └── main.go
│   ├── lesson23-subgraph-checkpoint/
│   │   └── main.go
│   ├── lesson24-graph-parallel-join/
│   │   └── main.go
│   ├── lesson25-extensible-scaffold/
│   │   └── main.go
│   ├── lesson26-chain-append-graph/
│   │   └── main.go
│   ├── lesson27-chain-parallel/
│   │   └── main.go
│   ├── lesson28-workflow-field-mapping/
│   │   └── main.go
│   ├── lesson29-tool-interrupt-resume/
│   │   └── main.go
│   ├── lesson30-tool-composite-interrupt/
│   │   └── main.go
│   ├── lesson31-custom-agent/
│   │   └── main.go
│   ├── lesson32-custom-resumable-agent/
│   │   └── main.go
│   ├── lesson33-custom-streaming-agent/
│   │   └── main.go
│   ├── lesson34-agent-session-values/
│   │   └── main.go
│   └── lesson35-agent-tool-advanced/
│       └── main.go
└── lessons/
    ├── lesson01-chatmodel.md
    ├── lesson02-stream.md
    ├── lesson03-prompt-template.md
    ├── lesson04-tool.md
    ├── lesson05-compose-chain.md
    ├── lesson06-tools-node.md
    ├── lesson07-model-with-tools.md
    ├── lesson08-workflow.md
    ├── lesson09-message-parser.md
    ├── lesson10-chatmodel-agent.md
    ├── lesson11-streamable-tool.md
    ├── lesson12-workflow-branch.md
    ├── lesson13-agent-stream.md
    ├── lesson14-agent-interrupt-resume.md
    ├── lesson15-workflow-checkpoint.md
    ├── lesson16-graph-basic.md
    ├── lesson17-graph-state.md
    ├── lesson18-agent-tool.md
    ├── lesson19-workflow-agents.md
    ├── lesson20-mini-project.md
    ├── lesson21-loop-agent.md
    ├── lesson22-graph-multibranch.md
    ├── lesson23-subgraph-checkpoint.md
    ├── lesson24-graph-parallel-join.md
    ├── lesson25-extensible-scaffold.md
    ├── lesson26-chain-append-graph.md
    ├── lesson27-chain-parallel.md
    ├── lesson28-workflow-field-mapping.md
    ├── lesson29-tool-interrupt-resume.md
    ├── lesson30-tool-composite-interrupt.md
    ├── lesson31-custom-agent.md
    ├── lesson32-custom-resumable-agent.md
    ├── lesson33-custom-streaming-agent.md
    ├── lesson34-agent-session-values.md
    └── lesson35-agent-tool-advanced.md
```

## 运行前提

- Go 1.18+
- 一个 OpenAI 兼容模型服务
- 环境变量：
  - `OPENAI_API_KEY`
  - `OPENAI_MODEL`
  - `OPENAI_BASE_URL`（可选；如果你不是直连 OpenAI，通常需要设置）

## 这一套教程的节奏

- 先学最小组件怎么单独使用
- 再学 Prompt、Stream、Tool、Agent、Compose
- 再往后进入 Graph、多 Agent、Checkpoint、Agent 实现细节和项目骨架
- 每一课只引入少量新概念，避免一次堆太多抽象
