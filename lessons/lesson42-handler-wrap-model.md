# 第 42 课：学会 WrapModel 和 SendEvent

这一课只讲一件事：让 handler 既能包裹 model，又能主动往事件流里发自定义事件。

学完这一课，你要记住：

```text
BeforeModelRewriteState -> SendEvent
WrapModel -> intercept Generate / Stream
```

## 1. 这一课解决什么问题

当你需要：

- 在模型调用前发可观测事件
- 给模型输出统一加包装
- 做 tracing / debug / telemetry

`WrapModel` 和 `SendEvent` 就很关键。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson42-handler-wrap-model/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson42-handler-wrap-model
```

## 4. 这节课最关键的代码

### 先用 `SendEvent(...)` 主动发事件

```go
event := adk.EventFromMessage(...)
err := adk.SendEvent(ctx, event)
```

这类事件会直接进入上层 `Runner` 的事件流。

### 再用 `WrapModel(...)` 包一层模型

```go
return &modelPrefixWrapper{inner: m}, nil
```

然后在 wrapper 里统一改输出：

```go
message.Content = "wrapped -> " + message.Content
```

## 5. 本课真正要记住的事

1. `SendEvent(...)` 适合中间态可观测事件
2. `WrapModel(...)` 适合统一包模型调用
3. 这两个点都只应该在 handler 里使用
4. 这是做 tracing、debug、审计、统一包装的关键接口
5. 如果只改 state，不需要上这两个接口

## 6. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
