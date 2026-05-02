# 第 28 课：学会 Workflow 字段映射

这一课只讲一件事：把 `FromField`、`MapFields`、`ToField` 放到一个最小 demo 里吃透。

学完这一课，你要记住：

```text
FromField / MapFields / ToField
```

## 1. 这一课解决什么问题

workflow 最容易让人绕住的地方，不是节点本身，而是：

- 上游结果到底怎么送给下游
- 是传整个对象
- 还是只传一个字段
- 还是把整个输出塞进某个目标字段

这节课就是专门把这三个映射动作拆开讲。

## 2. 本课 demo 在哪里

- 代码：`cmd/lesson28-workflow-field-mapping/main.go`

这节课不需要模型密钥。

## 3. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson28-workflow-field-mapping
```

## 4. 这节课最关键的代码

### 第一步：`FromField`

```go
AddInput(compose.START, compose.FromField("Topic"))
```

含义是：

- 从上游对象里拿出 `Topic`
- 作为当前节点的完整输入

### 第二步：`MapFields`

```go
AddInput(compose.START, compose.MapFields("Audience", "Value"))
```

含义是：

- 把上游的 `Audience`
- 映射到当前输入对象的 `Value`

### 第三步：`ToField`

```go
AddInput("normalize_topic", compose.ToField("Topic"))
AddInput(compose.START, compose.ToField("Original"))
```

含义是：

- 把整个上游输出
- 塞到当前输入对象的某个字段里

这一步很关键，因为很多人会把 `ToField` 和 `MapFields` 混掉。

## 5. 怎么区分这三个映射

`FromField`：

- 从上游拿一个字段
- 作为当前节点的完整输入

`MapFields`：

- 从上游拿一个字段
- 写到当前节点输入的某个字段

`ToField`：

- 拿上游整个输出
- 写到当前节点输入的某个字段

## 6. 本课真正要记住的事

1. workflow 的难点常常在字段映射，不在节点定义
2. `FromField` 是“字段 -> 整体输入”
3. `MapFields` 是“字段 -> 字段”
4. `ToField` 是“整体输出 -> 字段”
5. 一旦映射关系清楚，workflow 的可读性会明显提高

## 7. 官方资料

- Eino 仓库：https://github.com/cloudwego/eino
