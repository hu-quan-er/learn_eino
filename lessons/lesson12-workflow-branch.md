# 第 12 课：学会 Workflow 分支路由

这一课只讲一件事：让 workflow 不再是固定直线，而是能根据输入走不同分支。

学完这一课，你要记住：

```text
START -> Branch -> one selected node -> END
```

## 1. 这一课解决什么问题

第 08 课里的 workflow 还是静态编排：

- 节点关系提前写死
- 所有该跑的节点都会跑

但很多应用都需要路由：

- 命中某类问题就走 A 分支
- 否则走 B 分支

这就是 branch 的职责。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson12-workflow-branch/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson12-workflow-branch
```

你会看到两次调用：

- 一次命中 `answer_eino`
- 一次命中 `answer_general`

## 4. 这节课最关键的代码

### 第一步：先把候选分支节点都定义好

```go
workflow.AddLambdaNode("answer_eino", ...)
workflow.AddLambdaNode("answer_general", ...)
```

branch 本身不会替你创建节点，它只负责“选谁执行”。

### 第二步：定义 branch 条件

```go
workflow.AddBranch(compose.START, compose.NewGraphBranch(
    func(ctx context.Context, input string) (string, error) {
        if strings.Contains(strings.ToLower(input), "eino") {
            return "answer_eino", nil
        }
        return "answer_general", nil
    },
    map[string]bool{
        "answer_eino":    true,
        "answer_general": true,
    },
))
```

这段代码里你要记住两件事：

1. 条件函数返回的是“下一个节点名”
2. 允许跳转到哪些节点，要在 `map[string]bool` 里显式声明

### 第三步：把不同分支的输出都汇总到 END

```go
workflow.End().AddInput("answer_eino", compose.MapFields("answer", "answer"))
workflow.End().AddInput("answer_general", compose.MapFields("answer", "answer"))
```

branch 只负责选路。

最终要怎么把各分支结果收口，还是你自己定义。

## 5. 本课真正要记住的事

1. `AddBranch(...)` 是 workflow 的路由入口
2. `NewGraphBranch(...)` 的返回值是“目标节点 key”
3. branch 只会让命中的分支继续执行
4. branch 允许的目标节点必须提前声明
5. 分支后的结果仍然要自己汇总到 `END`

## 6. 什么时候该用 branch

适合：

- 意图路由
- 分类后分支处理
- 命中条件后提前结束

不适合：

- 只是单纯的字段映射
- 不存在不同执行路径的线性流程

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
