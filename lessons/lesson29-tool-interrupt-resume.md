# 第 29 课：学会 Tool Interrupt / Resume

这一课只讲一件事：让一个 tool 在执行时主动暂停，等待外部恢复后继续完成。

学完这一课，你要记住：

```text
tool.StatefulInterrupt -> checkpoint -> ResumeWithData
```

## 1. 这一课解决什么问题

很多真实工具都不是“一次调完就结束”，而是会遇到这种情况：

- 需要人工审批
- 需要等外部系统确认
- 需要补充一段恢复数据

这时候，tool 本身就应该能中断和恢复。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson29-tool-interrupt-resume/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson29-tool-interrupt-resume
```

你会看到两段输出：

1. 第一次运行时 tool 中断，并打印 interrupt ID / address
2. 恢复后 tool 返回最终结果

## 4. 这节课最关键的代码

### 第一步：第一次运行时保存状态并中断

```go
if !wasInterrupted {
    return "", tool.StatefulInterrupt(ctx, "need approval for "+callID, ApprovalState{Arguments: argumentsInJSON})
}
```

这一步会把 tool 自己的状态写进 checkpoint。

### 第二步：恢复时读取旧状态和恢复数据

```go
wasInterrupted, hasState, state := tool.GetInterruptState[ApprovalState](ctx)
isResumeTarget, hasData, decision := tool.GetResumeContext[string](ctx)
```

这里要分清两类信息：

- `GetInterruptState(...)`：上次中断时留下来的内部状态
- `GetResumeContext(...)`：这次恢复时外部喂进来的数据

### 第三步：只在自己是 resume target 时继续执行

```go
if !isResumeTarget {
    return "", tool.StatefulInterrupt(ctx, "waiting for resume "+callID, state)
}
```

这是多中断点场景里非常重要的习惯。

如果当前不是被恢复的那个点，就不要误继续，而是重新中断。

### 第四步：状态类型要能被 checkpoint 序列化

```go
schema.Register[ApprovalState]()
```

这点一定要记住。

只要你的 tool state 会进 checkpoint，就要提前注册可序列化类型。

## 5. 本课真正要记住的事

1. tool 也可以自己发起 interrupt
2. `StatefulInterrupt` 适合保存 tool 内部状态
3. `GetInterruptState` 和 `GetResumeContext` 不是一回事
4. 多中断点时要检查自己是不是 resume target
5. 进入 checkpoint 的自定义 state 需要先注册

## 6. 什么时候会用到它

- 审批类工具
- 支付、发布、删除等高风险动作
- 需要人工确认的外部系统调用

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
