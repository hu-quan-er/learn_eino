# 第 49 课：把 Interrupt / Resume 封装成 Service

这一课只讲一件事：不要把 checkpoint 和 resume 逻辑散落在调用方，应该包成 service。

学完这一课，你要记住：

```text
Service.StartPublish -> PendingApproval -> Service.ResumePublish
```

## 1. 这一课解决什么问题

如果每个调用方都自己处理：

- checkpoint id
- interrupt id
- resume params
- runner.ResumeWithParams(...)

代码会很快变乱。

更合理的方式是：

- 上层只调用业务 service
- service 内部管理 ADK 中断恢复细节

## 2. 本课 demo 在哪里

- 代码入口：`cmd/lesson49-interrupt-service/main.go`
- 内部包：`internal/lesson49service/`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson49-interrupt-service
```

## 4. 这节课最关键的结构

### service 对外暴露的是业务语义

```go
func (s *Service) StartPublish(ctx context.Context, request string) (*PendingApproval, error)
func (s *Service) ResumePublish(ctx context.Context, decision string) (string, error)
```

### 上层不直接碰 runner 细节

`main.go` 只关心：

- 有没有 pending approval
- 最终结果是什么

而不是自己去找 root cause interrupt。

## 5. 本课真正要记住的事

1. interrupt/resume 是基础设施细节，不该泄漏到每个调用方
2. service 层应该返回业务对象，比如 `PendingApproval`
3. checkpoint store 应该由 service 管
4. 这样后面接 API、接数据库、接人工审批系统都会更自然
5. 这比“教程式直接在 main 里 resume”更接近真实项目

## 6. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
