# 第 48 课：给 Agent 写单元测试

这一课只讲一件事：把 agent 逻辑放进可测试包里，并用 `go test` 验证。

学完这一课，你要记住：

```text
internal/package -> RunOnce helper -> go test
```

## 1. 这一课解决什么问题

如果 agent 逻辑只活在 `cmd/main.go` 里，你几乎没法认真测试。

想测试，就要先把逻辑放进：

- 可 import 的包
- 可重复调用的函数

## 2. 本课 demo 在哪里

- 代码入口：`cmd/lesson48-agent-testing/main.go`
- 被测包：`internal/lesson48testing/`
- 测试文件：`internal/lesson48testing/agent_test.go`

## 3. 如何运行

先跑 demo：

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson48-agent-testing
```

再跑测试：

```bash
go test ./internal/lesson48testing
```

## 4. 这节课最关键的结构

### 被测 agent 放在普通包里

```go
func NewReviewAgent(name, prefix string) *ReviewAgent
```

### 再给一个便于测试的 helper

```go
func RunOnce(ctx context.Context, name, prefix, query string) ([]string, error)
```

测试里直接调用它，而不是测试 `main.go`。

## 5. 本课真正要记住的事

1. `main.go` 不适合承载核心业务逻辑
2. 可测试代码应该放在普通包里
3. 给 agent 包提供一个小的执行 helper，测试会简单很多
4. 单元测试先测稳定逻辑，再测复杂编排
5. 工程化里，“能不能测”比“能不能跑”更重要

## 6. 一个实际验证结果

这节课的 `go test ./internal/lesson48testing` 我已经跑过，能通过。

说明一下环境现象：

- 沙箱里的 `go test` 报过一次 `testing/internal/testdeps: package testmain: cannot find package`
- 换到非沙箱执行后通过

所以这个失败不是课件代码问题，而是当前测试执行环境差异。

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
