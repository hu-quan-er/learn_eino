# 第 47 课：学会用 Runner Factory 做依赖注入

这一课只讲一件事：把“怎么造 runner”抽成一个工厂，而不是每处手写一遍。

学完这一课，你要记住：

```text
Dependencies -> RunnerFactory -> Service -> Query
```

## 1. 这一课解决什么问题

真实项目里，runner 构造通常依赖：

- 默认 session values
- 统一 checkpoint store
- 公共 callbacks
- 通用依赖对象

这些不应该散落在每个调用点里。

## 2. 本课 demo 在哪里

- 代码入口：`cmd/lesson47-runner-factory/main.go`
- 内部包：`internal/lesson47factory/`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson47-runner-factory
```

## 4. 这节课最关键的结构

### 依赖对象

```go
type Dependencies struct {
    SessionValues map[string]any
}
```

### 工厂

```go
type RunnerFactory struct {
    deps Dependencies
}
```

### 服务对象

```go
type Service struct {
    runner *adk.Runner
    opts   []adk.AgentRunOption
}
```

工厂负责统一装配，服务对象负责执行业务查询。

## 5. 本课真正要记住的事

1. runner factory 适合收拢公共依赖
2. 默认 `AgentRunOption` 也可以作为依赖注入进去
3. 同一个工厂可以生成多个不同配置的 service
4. 这样能减少复制粘贴和配置漂移
5. 这是比“每次直接 `adk.NewRunner(...)`”更像工程代码的写法

## 6. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
