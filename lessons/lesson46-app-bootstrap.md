# 第 46 课：做一个 App Bootstrap 分层

这一课只讲一件事：不要让 `main.go` 直接堆满 agent 构造细节，先学会做一个最小 bootstrap 层。

学完这一课，你要记住：

```text
cmd/main -> internal/app.New -> runner -> app.Run
```

## 1. 这一课解决什么问题

教程里单文件 demo 很适合学 API。

但到了真实项目里，如果继续把：

- config
- agent 构造
- runner 创建
- 运行逻辑

全塞进 `main.go`，很快就会失控。

## 2. 本课 demo 在哪里

- 代码入口：`cmd/lesson46-app-bootstrap/main.go`
- 内部包：`internal/lesson46app/`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson46-app-bootstrap
```

## 4. 这节课最关键的结构

### `Config`

```go
type Config struct {
    AgentName string
    Prefix    string
}
```

### `New(...)`

```go
func New(ctx context.Context, cfg Config) *App
```

### `Run(...)`

```go
func (a *App) Run(ctx context.Context, query string) ([]string, error)
```

这三个点说明：

- `main.go` 只负责启动
- 内部包负责装配
- 业务运行入口是 `App.Run(...)`

## 5. 本课真正要记住的事

1. bootstrap 的第一目标是把装配逻辑从 `main.go` 挪出去
2. `cmd` 目录只保留启动代码
3. `internal/...` 放应用实现
4. 先做这一步，后面加 config、service、测试都会轻很多
5. 工程化不是“多写文件”，而是让职责边界清晰

## 6. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
