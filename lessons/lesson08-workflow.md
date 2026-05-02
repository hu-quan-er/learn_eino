# 第 08 课：学会 Workflow 基础编排

这一课只讲一件事：理解 `Workflow` 和 `Chain` 的关系。

学完这一课，你要记住：

```text
START -> nodes with dependencies -> END
```

## 1. 为什么第 05 课有 Chain，这一课还要学 Workflow

`Chain` 适合线性顺序：

- A -> B -> C

但真实编排很快会遇到：

- 一个输入拆给多个节点
- 多个节点结果汇总到最终输出

这时 `Workflow` 更自然，因为它本质上是一个带依赖关系的图。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson08-workflow/main.go`

这节课是纯本地 demo，不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson08-workflow
```

你会看到类似输出：

```text
workflow output: map[string]interface {}{"normalized_question":"什么是 Eino Workflow？", "note":"Alice is learning Eino workflow."}
```

## 4. 这节课的结构怎么理解

这一课的 workflow 很简单：

1. `normalize_question` 节点处理问题文本
2. `build_note` 节点处理学生信息
3. `END` 汇总两个节点的结果

也就是说，这不是线性链，而是一个很小的扇出再汇合的图。

## 5. 这节课最关键的代码

### 第一步：创建 Workflow

```go
workflow := compose.NewWorkflow[map[string]any, map[string]any]()
```

这里定义了整个图的输入输出：

- 输入是 `map[string]any`
- 输出也是 `map[string]any`

### 第二步：定义节点

```go
workflow.AddLambdaNode("normalize_question", ...)
workflow.AddLambdaNode("build_note", ...)
```

这里仍然用的是 `Lambda`，因为我们只想先聚焦在“图的连接方式”，而不是节点实现本身。

### 第三步：声明输入映射

```go
AddInput(compose.START, compose.MapFields("question", "Question"))
```

这一行很关键。它的意思是：

- 从 `START` 取输入字段 `question`
- 映射到当前节点输入结构里的 `Question`

也就是说，`Workflow` 不只是表达执行顺序，还表达数据怎么流动。

### 第四步：汇总到 END

```go
workflow.End().AddInput("normalize_question", compose.MapFields("normalized_question", "normalized_question"))
workflow.End().AddInput("build_note", compose.MapFields("note", "note"))
```

这一步表示：

- 把多个节点的输出汇总到最终输出

`END` 可以理解成 workflow 的最终结果收集点。

## 6. 本课真正要记住的事

只记这五点：

1. `Workflow` 是图式编排，不只是顺序编排
2. `START` 是输入入口
3. `END` 是输出汇总点
4. `AddInput(...)` 同时表达依赖关系和数据映射
5. `MapFields(...)` 用来说明字段怎么从上游映射到下游

## 7. 下一课讲什么

下一课讲 `MessageParser`。

因为到这里你已经把消息生成和流程编排串起来了，下一步该学会“怎么把模型返回的文本变成结构化数据”。

## 8. 官方资料

- Chain & Graph Orchestration: https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/
- Eino 仓库: https://github.com/cloudwego/eino
