# 第 05 课：用 Chain 串起 Prompt 和 Model

这一课只讲一件事：开始进入 Eino 的编排层。

学完这一课，你要记住这条最小链路：

```text
input vars -> ChatTemplate -> ChatModel -> Lambda -> output
```

## 1. 这一课解决什么问题

前几课你学的是单个组件：

- `ChatModel`
- `Stream`
- `Prompt`
- `Tool`

但真实应用不是单个组件独立工作，而是多个组件顺序衔接。

这时候你当然可以自己手写：

1. 先调 Prompt
2. 再把结果交给 Model
3. 再手动提取输出

但代码会越来越碎。

所以 Eino 提供了 `compose` 层。对最常见的线性流程，最简单的入口就是 `Chain`。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson05-compose-chain/main.go`

这个 demo 会把三件事连起来：

1. `ChatTemplate`
2. `ChatModel`
3. 一个把 `*schema.Message` 转成 `string` 的 `Lambda`

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
export OPENAI_API_KEY="你的密钥"
export OPENAI_MODEL="gpt-4o-mini"
export OPENAI_BASE_URL="你的兼容接口地址"
go run ./cmd/lesson05-compose-chain
```

如果配置正常，你会看到类似输出：

```text
chain output:
1. ...
2. ...
3. ...
```

## 4. 这节课最关键的认识

`Chain` 不是新模型，也不是新 Prompt。

你可以把它理解成：

- 一个顺序编排容器
- 负责把上一个节点的输出交给下一个节点

所以这节课的核心不是某个单独 API，而是“组件开始被串起来了”。

## 5. 代码怎么读

### 第一步：准备 Prompt

```go
chatTemplate := prompt.FromMessages(
    schema.FString,
    schema.SystemMessage("你是一个{role}，回答时给出 3 条要点。"),
    schema.UserMessage("{question}"),
)
```

这部分你在第三课已经见过了。

它的职责不变：

- 输入变量
- 输出 `[]*schema.Message`

### 第二步：创建 Chain

```go
chain := compose.NewChain[map[string]any, string]()
```

这里的泛型非常关键：

- 输入是 `map[string]any`
- 输出是 `string`

这就是整个链路的“对外接口”。

### 第三步：顺序追加节点

```go
chain.
    AppendChatTemplate(chatTemplate).
    AppendChatModel(chatModel).
    AppendLambda(compose.InvokableLambda(func(ctx context.Context, message *schema.Message) (string, error) {
        return message.Content, nil
    }))
```

这是本课最核心的代码块。

它表达的是一个非常清晰的执行顺序：

1. 先把变量格式化成消息
2. 再把消息交给模型
3. 再把模型返回的 `*schema.Message` 转成纯字符串

最后这一步用的是 `Lambda`，你现在可以先把它理解成“自定义转换节点”。

### 第四步：Compile

```go
runner, err := chain.Compile(ctx)
```

`Compile(...)` 的意思不是“编译 Go 代码”，而是：

- 把你声明的链路结构整理成一个可运行对象

你可以把 `runner` 理解成“已经装配好的执行器”。

### 第五步：Invoke

```go
answer, err := runner.Invoke(ctx, map[string]any{
    "role":     "Eino 助教",
    "question": "请简要解释为什么 Eino 要把 Prompt 和 Model 拆成不同组件。",
})
```

这时候整个链条就跑起来了：

- `role` 和 `question` 先进入 Prompt
- Prompt 输出消息
- Model 生成回答
- Lambda 抽取 `message.Content`
- 最终得到 `string`

## 6. 为什么这一课很重要

到第五课为止，你第一次真正看到 Eino 的“框架味道”。

因为前面几课都还是单个组件学习，而这一课开始体现 Eino 的核心价值：

- 组件标准化
- 组件可编排
- 上下游接口清晰

这也是为什么 Eino 不想让你把 Prompt、Model、Tool 全都揉在一个函数里。

## 7. 本课真正要记住的事

只记这五点：

1. `Chain` 用来表达顺序编排
2. `AppendChatTemplate`、`AppendChatModel`、`AppendLambda` 就是在按顺序搭流水线
3. `Compile` 会得到一个可执行的 runner
4. `Invoke` 才是真正执行整条链
5. `Lambda` 很适合做轻量的输入输出转换

## 8. 你现在可以自己做的练习

建议你立刻改三处：

1. 在 `ChatTemplate` 里增加一个 `{style}` 变量
2. 把最后的 `Lambda` 改成返回 `*schema.Message`，观察链的输出类型如何变化
3. 把用户问题换成 “Prompt、Tool、Compose 分别负责什么”

## 9. 到第 5 课为止，你已经掌握了什么

现在你已经有了一个比较完整的 Eino 入门骨架：

1. 第 01 课：直接调用 `ChatModel`
2. 第 02 课：流式读取 `Stream`
3. 第 03 课：用 `Prompt` 组织输入
4. 第 04 课：定义 `Tool`
5. 第 05 课：用 `Chain` 串起多个组件

这时你再继续往后学 `ToolsNode` 或 `Agent`，会顺很多。

## 10. 官方资料

本课内容参考的是 Eino 官方一手资料和当前版本源码：

- Chain & Graph Orchestration: https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/
- ChatTemplate Guide: https://www.cloudwego.io/docs/eino/core_modules/components/chat_template_guide/
- ChatModel Guide: https://www.cloudwego.io/docs/eino/core_modules/components/chat_model_guide/
- Eino 仓库: https://github.com/cloudwego/eino
