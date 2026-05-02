# 第 43 课：学会 Run Local Values

这一课只讲一件事：理解 `SetRunLocalValue / GetRunLocalValue / DeleteRunLocalValue` 这组运行期局部状态接口。

学完这一课，你要记住：

```text
SetRunLocalValue -> GetRunLocalValue -> DeleteRunLocalValue
```

## 1. 这一课解决什么问题

有些数据：

- 不是全局 session
- 也不是 message history
- 只是当前一次 agent 运行过程中的局部状态

这种数据就适合放在 run-local 里。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson43-run-local-values/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson43-run-local-values
```

## 4. 这节课最关键的代码

### 设置值

```go
adk.SetRunLocalValue(ctx, "trace_id", "trace-lesson43")
```

### 读取值

```go
value, found, err := adk.GetRunLocalValue(ctx, "trace_id")
```

### 删除值

```go
adk.DeleteRunLocalValue(ctx, "trace_id")
```

本课在 `AfterModelRewriteState(...)` 里把 run-local 的读取结果写进 session，再由下一个 agent 读出来，所以你能直接看到它确实存在并且被删掉了。

## 5. 本课真正要记住的事

1. run-local value 是“当前这次运行”的局部状态
2. 它和 session value 不是一回事
3. 这组 API 只能在 ChatModelAgent handler 执行过程中调用
4. 它适合 trace id、临时标记、局部缓存这类数据
5. 不需要跨 agent 共享时，优先考虑 run-local 而不是 session

## 6. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
