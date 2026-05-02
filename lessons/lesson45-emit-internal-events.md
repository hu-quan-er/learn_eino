# 第 45 课：学会 EmitInternalEvents

这一课只讲一件事：控制嵌套 `AgentTool` 的内部事件要不要透给最外层调用方。

学完这一课，你要记住：

```text
AgentTool + ChatModelAgent + ToolsConfig.EmitInternalEvents
```

## 1. 这一课解决什么问题

当一个 agent 调另一个 agent 时，你通常会遇到一个选择：

- 只看最外层最终结果
- 还是把内部 agent 的事件也实时透出来

`EmitInternalEvents` 就是这个开关。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson45-emit-internal-events/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson45-emit-internal-events
```

你会看到两组结果：

1. `EmitInternalEvents=false`
2. `EmitInternalEvents=true`

## 4. 这节课最关键的代码

### 关闭内部事件透传

```go
ToolsConfig: adk.ToolsConfig{
    EmitInternalEvents: false,
}
```

### 开启内部事件透传

```go
ToolsConfig: adk.ToolsConfig{
    EmitInternalEvents: true,
}
```

其他代码完全一样，只有这个开关不同。

## 5. 本课真正要记住的事

1. `EmitInternalEvents` 只影响内部 agent 事件是否向外透出
2. 它不改变工具调用逻辑本身
3. 调试多 agent 时，通常应该先开
4. 正式面向终端用户时，要根据产品需求决定是否保留
5. 这对流式 UI、调试界面、trace 页面都很关键

## 6. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
