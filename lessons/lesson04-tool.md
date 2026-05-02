# 第 04 课：学会定义一个 Tool

这一课只讲一件事：理解 Eino 里的 Tool 到底由什么组成。

学完这一课，你要记住：

```text
Tool = ToolInfo(schema) + InvokableRun(execution)
```

## 1. 这一课解决什么问题

模型本身只能“生成文本”，但真实应用经常需要：

- 查天气
- 查数据库
- 调内部服务
- 调本地函数

这类“模型之外的能力”，在 Eino 里就叫 `Tool`。

所以你现在要先建立一个正确认识：

- Tool 不是提示词
- Tool 也不是 Agent
- Tool 是一种“可被模型调用的外部能力”

## 2. Tool 为什么至少有两部分

一个工具想被模型调用，最少要有两件东西：

1. 模型得知道这个工具叫什么、接收什么参数
2. 系统得真的能把它执行起来

在 Eino 里，这两部分分别对应：

- `Info()`：返回 `schema.ToolInfo`
- `InvokableRun(...)`：真正执行工具

所以这节课的重点不是“让模型决定调工具”，而是先学会“把工具定义清楚”。

## 3. 本课 demo 在哪里

- 代码：`cmd/lesson04-tool/main.go`

这一课也是本地 demo，不依赖模型密钥。

## 4. 如何运行

```bash
cd /Users/seaya/huquan.95/project-4/eino
go run ./cmd/lesson04-tool
```

你会看到两类输出：

1. 工具推断出来的输入 schema
2. 工具执行后的结果 JSON

类似这样：

```text
tool name: get_weather
tool desc: return mock weather data by city

inferred input schema:
{
  ...
}

tool result:
{"summary":"Shanghai is sunny today","temperature":26,"unit":"celsius"}
```

## 5. 这一课最关键的代码

### 第一步：定义输入输出结构

```go
type WeatherQuery struct {
    City string `json:"city" jsonschema:"required,description=city name"`
    Unit string `json:"unit" jsonschema:"description=temperature unit,enum=celsius,enum=fahrenheit"`
}
```

这里定义的是“工具的参数结构”。

Eino 的 `InferTool(...)` 会根据这个结构自动推断参数 schema。

所以这一步很重要，因为它决定了模型将来看到的工具参数长什么样。

### 第二步：写工具逻辑

```go
func getWeather(ctx context.Context, input *WeatherQuery) (*WeatherResult, error) {
    ...
}
```

这就是工具真正的业务逻辑。

你现在可以先把它理解成一个普通 Go 函数，只是它后面会被包装成 Eino Tool。

### 第三步：把函数包装成 Tool

```go
weatherTool, err := toolutils.InferTool(
    "get_weather",
    "return mock weather data by city",
    getWeather,
)
```

这一行是本课最关键的一行。

`InferTool(...)` 做了两件事：

1. 根据 `WeatherQuery` 自动生成工具参数 schema
2. 把 `getWeather` 包装成一个可执行的 `InvokableTool`

也就是说，你不用手写一大堆 `ToolInfo` 样板代码。

### 第四步：查看工具元信息

```go
info, err := weatherTool.Info(ctx)
schemaJSON, err := info.ToJSONSchema()
```

这一步是在看 Tool 的“对外描述”。

这份 schema 将来就是模型理解工具参数的依据。

### 第五步：直接执行工具

```go
result, err := weatherTool.InvokableRun(
    ctx,
    `{"city":"Shanghai","unit":"celsius"}`,
)
```

这里你要抓住一个很重要的事实：

- `InvokableRun(...)` 收到的是 JSON 字符串
- 返回的也是字符串结果

这正是标准 `InvokableTool` 的特点。

## 6. 本课真正要记住的事

只记这五点：

1. Tool 是模型之外的外部能力
2. Tool 至少包含“schema 描述”和“执行逻辑”
3. `Info()` 负责告诉系统这个工具怎么用
4. `InvokableRun(...)` 负责真正执行工具
5. `InferTool(...)` 能从 Go 结构体自动推断参数 schema

## 7. 这一课和 Agent 有什么关系

先不要急着把 Tool 和 Agent 混在一起。

你现在只需要先理解：

- Tool 是单个能力
- ToolsNode 是工具执行器
- Agent 是会决定“什么时候调哪个工具”的更高层结构

所以这一课只是先把“能力本身”学会。

## 8. 你现在可以自己做的练习

建议你立刻改三处：

1. 给 `WeatherQuery` 新增 `Date string`
2. 把 `Unit` 改成只允许 `celsius`
3. 再定义一个 `get_exchange_rate` 工具，自己试着用 `InferTool` 包起来

## 9. 下一课讲什么

下一课讲 `Compose Chain`。

因为到这一步你已经有了：

- Prompt
- Model
- Tool 的基础概念

下一步该学“怎么把组件顺序串起来”，也就是 Eino 的编排层。

## 10. 官方资料

本课内容参考的是 Eino 官方一手资料和当前版本源码：

- Tools Guide: https://www.cloudwego.io/docs/eino/core_modules/components/tools_node_guide/
- How to Create a Tool: https://www.cloudwego.io/docs/eino/core_modules/components/tools_node_guide/how_to_create_a_tool/
- Eino 仓库: https://github.com/cloudwego/eino
