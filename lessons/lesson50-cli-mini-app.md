# 第 50 课：做一个完整 CLI 小项目

这一课只讲一件事：把前面的 agent、session、分层整理成一个最小可运行 CLI 小应用。

学完这一课，你要记住：

```text
cmd -> internal/app -> workflow -> session values -> final result
```

## 1. 这一课解决什么问题

到第 50 课为止，你已经学了很多点状能力。

这一课要把它们收起来，变成一个最小完整体：

- 有 `cmd` 入口
- 有 `internal/app`
- 有 workflow
- 有 session values
- 有可打印的 trace

## 2. 本课 demo 在哪里

- 代码入口：`cmd/lesson50-cli-mini-app/main.go`
- 内部包：`internal/lesson50app/`

这节课不需要模型密钥。

## 3. 如何运行

默认运行：

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson50-cli-mini-app
```

也可以自己带 query：

```bash
go run ./cmd/lesson50-cli-mini-app 请生成一个 agent 工程化总结
```

## 4. 这节课最关键的结构

### `Config`

```go
type Config struct {
    AppName string
    Tenant  string
}
```

### `App`

```go
func New(ctx context.Context, cfg Config) (*App, error)
func (a *App) Run(ctx context.Context, query string) (*Result, error)
```

### workflow agents

- `plannerAgent`
- `writerAgent`

planner 先把 plan 写进 session，writer 再读出来生成最终结果。

## 5. 本课真正要记住的事

1. 小项目不需要一开始就很复杂，但结构必须清楚
2. `cmd`、`internal/app`、agents、result model 这些边界值得尽早固定
3. session values 很适合承载跨 agent 的内部上下文
4. trace 收集不要等到系统很大才做
5. 第 50 课的价值不在功能多，而在结构已经像一个真正项目

## 6. 学完第 50 课后该怎么继续

到这里，继续往下再拆“第 51 课、第 52 课”当然还能写。

但更高价值的方向通常变成：

- 专题项目
- 某一类 middleware 深挖
- 某一个真实 agent 系统重构

也就是说，后面更适合从“章节教程”转到“专题实战”。

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
