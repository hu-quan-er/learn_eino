# 第 03 课：学会 Prompt Template

这一课只讲一件事：把“写死的消息”变成“可复用的消息模板”。

学完这一课，你只要记住这一条链路：

```text
variables -> ChatTemplate -> []*schema.Message
```

## 1. 这一课解决什么问题

前两课里，你的消息基本都是直接写在代码里的：

- system 写死
- user 写死
- 每次改问题都要改代码

这在 demo 阶段没问题，但一旦你开始做真实应用，会立刻遇到两个问题：

1. 输入不稳定，变量散落在代码里
2. 同一套提示结构没法复用

所以 Eino 把 Prompt 抽象成 `ChatTemplate`。

它的职责非常单纯：

- 接收变量
- 填充模板
- 产出标准 `[]*schema.Message`

注意这里先不要把 Prompt 想成“生成答案”。它不负责调用模型，它只负责把输入整理好。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson03-prompt-template/main.go`

这一课是一个纯本地 demo，不需要模型密钥，也不需要联网调用。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson03-prompt-template
```

你会看到类似输出：

```text
formatted messages:
1. role=system
   content=你是一个Eino 入门老师。回答要短、准、清晰。

2. role=user
   content=我已经学完了 Generate 和 Stream。

3. role=assistant
   content=很好，下一步应该学会怎么把输入组织成模板。

4. role=user
   content=请基于上面的上下文回答：Prompt Template 在 Eino 里负责什么？
```

这个输出很关键，因为它说明：

- Prompt 的输出不是字符串
- 而是一组已经排好顺序的 `schema.Message`

## 4. 最核心的代码是什么

### 第一步：创建模板

```go
chatTemplate := prompt.FromMessages(
    schema.FString,
    schema.SystemMessage("你是一个{role}。回答要短、准、清晰。"),
    schema.MessagesPlaceholder("history", true),
    schema.UserMessage("请基于上面的上下文回答：{question}"),
)
```

这一段就是本课核心。

你可以把它理解成：

- `prompt.FromMessages(...)`：把多条消息拼成一个模板
- `schema.FString`：告诉 Eino 用 `{variable}` 这种格式做变量替换
- `schema.MessagesPlaceholder(...)`：在固定位置插入一段消息历史

这里最值得注意的是 `MessagesPlaceholder("history", true)`。

它的意思是：

- 模板里预留一个叫 `history` 的位置
- 后续可以把 `[]*schema.Message` 塞进来
- `true` 表示这个位置是可选的

## 5. Format 到底做了什么

```go
messages, err := chatTemplate.Format(ctx, map[string]any{
    "role": "Eino 入门老师",
    "history": []*schema.Message{
        schema.UserMessage("我已经学完了 Generate 和 Stream。"),
        {
            Role:    schema.Assistant,
            Content: "很好，下一步应该学会怎么把输入组织成模板。",
        },
    },
    "question": "Prompt Template 在 Eino 里负责什么？",
})
```

`Format(...)` 的作用很直接：

- 把变量 map 填到模板里
- 返回格式化后的 `[]*schema.Message`

所以你现在应该把 Prompt 理解成“消息构造器”，不是“模型调用器”。

## 6. 这一课最容易踩的坑

最常见的坑只有两个：

1. 模板里用了 `{question}`，但变量 map 里没给 `question`
2. `MessagesPlaceholder` 需要的是 `[]*schema.Message`，你却传了别的类型

也就是说，Prompt 的报错很多时候不是模型报错，而是你自己的输入结构不匹配。

## 7. 本课真正要记住的事

只记这四点：

1. `ChatTemplate` 不负责调用模型
2. `Format` 的输出是 `[]*schema.Message`
3. `schema.FString` 用 `{var}` 方式替换变量
4. `schema.MessagesPlaceholder` 用来插入消息历史

## 8. 你现在可以自己做的练习

建议你立刻改三处：

1. 把 `role` 改成 “严谨的 Go 架构师”
2. 把 `history` 去掉，观察可选 placeholder 的行为
3. 故意删掉 `question`，看看报错长什么样

## 9. 下一课讲什么

下一课讲 `Tool`。

因为到这一步，你已经会了：

- 组织输入
- 调用模型
- 流式读取

下一步就该理解“模型之外的能力怎么被接进 Eino”。

## 10. 官方资料

本课内容参考的是 Eino 官方一手资料和当前版本源码：

- ChatTemplate Guide: https://www.cloudwego.io/docs/eino/core_modules/components/chat_template_guide/
- Components Overview: https://www.cloudwego.io/docs/eino/core_modules/components/
- Eino 仓库: https://github.com/cloudwego/eino
