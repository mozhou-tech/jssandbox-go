# Eino 大模型编排规则

## 概述

本项目使用 [CloudWeGo Eino](https://github.com/cloudwego/eino) 作为大模型编排框架。Eino 是字节跳动开源的 Golang 大模型应用开发框架，提供组件化设计、图编排引擎和流式处理机制。

## 核心概念

### 组件化设计
- Eino 将常见功能模块抽象为组件（如 ChatModel、Lambda 等）
- 每个组件有多种实现，支持嵌套和复杂业务逻辑构建
- 组件可以组合使用，形成复杂的工作流

### 图编排引擎
- 基于有向图描述组件之间的流转关系
- 支持复杂的业务流程编排
- 可以直观地构建和管理业务逻辑

### 流式处理
- 支持流式输入和输出
- 自动处理流的拼接和转换
- 提升应用的实时性和性能

## 依赖包

项目已添加以下 eino 相关依赖：

```go
require (
    github.com/cloudwego/eino v0.7.0
    github.com/cloudwego/eino-ext/callbacks/cozeloop v0.1.6
    github.com/cloudwego/eino-ext/components/document/parser/html v0.0.0-20251117090452-bd6375a0b3cf
    github.com/cloudwego/eino-ext/components/document/parser/pdf v0.0.0-20251117090452-bd6375a0b3cf
)
```

## 代码组织

### 目录结构

```
packages/sidecar/
├── handlers/
│   └── llm/              # LLM 相关处理器
│       ├── orchestration.go  # 编排逻辑
│       └── workflows.go      # 工作流定义
├── models/
│   └── llm/              # LLM 相关模型
│       ├── workflow.go   # 工作流模型
│       └── component.go  # 组件模型
└── services/
    └── llm/              # LLM 服务
        ├── eino.go       # Eino 初始化
        └── chat.go       # 聊天服务
```

## 基本使用模式

### 1. 初始化 Eino 引擎

```go
package services

import (
    "context"
    "github.com/cloudwego/eino"
    "github.com/cloudwego/eino-ext/callbacks/cozeloop"
)

// InitEinoEngine 初始化 Eino 引擎
func InitEinoEngine() (*eino.Engine, error) {
    engine := eino.NewEngine()
    
    // 注册回调（可选）
    engine.RegisterCallback(cozeloop.NewCallback())
    
    return engine, nil
}
```

### 2. 创建简单的工作流

```go
package handlers

import (
    "context"
    "github.com/cloudwego/eino"
    "github.com/cloudwego/eino/components/model"
)

// SimpleChatWorkflow 创建简单的聊天工作流
func SimpleChatWorkflow(engine *eino.Engine) (*eino.Workflow, error) {
    workflow := engine.NewWorkflow("simple_chat")
    
    // 创建聊天模型组件
    chatModel, err := model.NewChatModel(&model.ChatModelConfig{
        ModelName: "gpt-4",
        // 其他配置...
    })
    if err != nil {
        return nil, err
    }
    
    // 添加节点
    inputNode := workflow.AddNode("input", eino.NewInputNode())
    chatNode := workflow.AddNode("chat", chatModel)
    outputNode := workflow.AddNode("output", eino.NewOutputNode())
    
    // 连接节点
    workflow.Connect(inputNode, chatNode)
    workflow.Connect(chatNode, outputNode)
    
    return workflow, nil
}
```

### 3. 执行工作流

```go
// ExecuteWorkflow 执行工作流
func ExecuteWorkflow(ctx context.Context, workflow *eino.Workflow, input interface{}) (interface{}, error) {
    // 创建执行上下文
    execCtx := eino.NewExecutionContext(ctx)
    
    // 设置输入
    execCtx.SetInput(input)
    
    // 执行工作流
    result, err := workflow.Execute(execCtx)
    if err != nil {
        return nil, err
    }
    
    return result, nil
}
```

### 4. 文档解析工作流

```go
package handlers

import (
    "github.com/cloudwego/eino"
    "github.com/cloudwego/eino-ext/components/document/parser/pdf"
    "github.com/cloudwego/eino-ext/components/document/parser/html"
)

// DocumentProcessingWorkflow 创建文档处理工作流
func DocumentProcessingWorkflow(engine *eino.Engine) (*eino.Workflow, error) {
    workflow := engine.NewWorkflow("document_processing")
    
    // PDF 解析器
    pdfParser, err := pdf.NewParser()
    if err != nil {
        return nil, err
    }
    
    // HTML 解析器
    htmlParser, err := html.NewParser()
    if err != nil {
        return nil, err
    }
    
    // 添加节点
    inputNode := workflow.AddNode("input", eino.NewInputNode())
    pdfNode := workflow.AddNode("pdf_parser", pdfParser)
    htmlNode := workflow.AddNode("html_parser", htmlParser)
    mergeNode := workflow.AddNode("merge", eino.NewLambdaNode(func(ctx context.Context, input interface{}) (interface{}, error) {
        // 合并解析结果
        return mergeResults(input), nil
    }))
    outputNode := workflow.AddNode("output", eino.NewOutputNode())
    
    // 连接节点（根据文档类型路由）
    workflow.Connect(inputNode, pdfNode)
    workflow.Connect(inputNode, htmlNode)
    workflow.Connect(pdfNode, mergeNode)
    workflow.Connect(htmlNode, mergeNode)
    workflow.Connect(mergeNode, outputNode)
    
    return workflow, nil
}
```

## 流式处理

### 流式输出

```go
// StreamChatWorkflow 创建流式聊天工作流
func StreamChatWorkflow(engine *eino.Engine) (*eino.Workflow, error) {
    workflow := engine.NewWorkflow("stream_chat")
    
    chatModel, err := model.NewChatModel(&model.ChatModelConfig{
        ModelName: "gpt-4",
        Stream: true, // 启用流式输出
    })
    if err != nil {
        return nil, err
    }
    
    inputNode := workflow.AddNode("input", eino.NewInputNode())
    chatNode := workflow.AddNode("chat", chatModel)
    outputNode := workflow.AddNode("output", eino.NewOutputNode())
    
    workflow.Connect(inputNode, chatNode)
    workflow.Connect(chatNode, outputNode)
    
    return workflow, nil
}

// 流式执行
func StreamExecute(ctx context.Context, workflow *eino.Workflow, input interface{}, streamChan chan<- string) error {
    execCtx := eino.NewExecutionContext(ctx)
    execCtx.SetInput(input)
    
    // 注册流式回调
    execCtx.OnStream(func(chunk string) {
        streamChan <- chunk
    })
    
    _, err := workflow.Execute(execCtx)
    close(streamChan)
    return err
}
```

## 错误处理

### 工作流错误处理

```go
// ExecuteWithErrorHandling 带错误处理的工作流执行
func ExecuteWithErrorHandling(ctx context.Context, workflow *eino.Workflow, input interface{}) (interface{}, error) {
    execCtx := eino.NewExecutionContext(ctx)
    execCtx.SetInput(input)
    
    // 设置错误处理回调
    execCtx.OnError(func(err error) {
        log.Printf("工作流执行错误: %v", err)
        // 可以在这里实现重试、降级等逻辑
    })
    
    result, err := workflow.Execute(execCtx)
    if err != nil {
        return nil, fmt.Errorf("工作流执行失败: %w", err)
    }
    
    return result, nil
}
```

## 配置管理

### Eino 配置

在 `config/config.go` 中添加 Eino 相关配置：

```go
type Config struct {
    // ... 其他配置
    
    // Eino 配置
    Eino struct {
        // 模型配置
        DefaultModel string `yaml:"default_model"`
        ModelAPIKey  string `yaml:"model_api_key"`
        
        // 工作流配置
        WorkflowTimeout int `yaml:"workflow_timeout"` // 超时时间（秒）
        MaxConcurrency   int `yaml:"max_concurrency"`  // 最大并发数
        
        // 流式配置
        StreamBufferSize int `yaml:"stream_buffer_size"`
    } `yaml:"eino"`
}
```

## 最佳实践

### 1. 工作流设计

- **单一职责**: 每个工作流应该专注于一个特定的任务
- **可复用性**: 设计可复用的组件和工作流
- **错误处理**: 在工作流中适当位置添加错误处理节点
- **性能优化**: 对于可以并行的操作，使用并行节点

### 2. 组件使用

- **组件选择**: 根据需求选择合适的组件实现
- **配置管理**: 将组件配置外部化，便于管理
- **资源清理**: 确保组件资源正确释放

### 3. 流式处理

- **缓冲区管理**: 合理设置流式缓冲区大小
- **错误恢复**: 实现流式处理中的错误恢复机制
- **性能监控**: 监控流式处理的性能指标

### 4. 日志和监控

```go
// 在工作流中添加日志节点
logNode := workflow.AddNode("log", eino.NewLambdaNode(func(ctx context.Context, input interface{}) (interface{}, error) {
    logger.WithFields(logrus.Fields{
        "workflow": workflow.Name(),
        "input": input,
    }).Info("工作流执行")
    return input, nil
}))
```

## 示例：招标文档分析工作流

```go
// TenderDocumentAnalysisWorkflow 招标文档分析工作流
func TenderDocumentAnalysisWorkflow(engine *eino.Engine) (*eino.Workflow, error) {
    workflow := engine.NewWorkflow("tender_document_analysis")
    
    // 1. 文档解析节点
    pdfParser, _ := pdf.NewParser()
    parseNode := workflow.AddNode("parse", pdfParser)
    
    // 2. 文档提取节点（提取关键信息）
    extractNode := workflow.AddNode("extract", eino.NewLambdaNode(func(ctx context.Context, input interface{}) (interface{}, error) {
        // 提取招标要求、截止日期等关键信息
        return extractTenderInfo(input), nil
    }))
    
    // 3. 分析节点（使用 LLM 分析）
    chatModel, _ := model.NewChatModel(&model.ChatModelConfig{
        ModelName: "gpt-4",
    })
    analysisNode := workflow.AddNode("analyze", chatModel)
    
    // 4. 格式化输出节点
    formatNode := workflow.AddNode("format", eino.NewLambdaNode(func(ctx context.Context, input interface{}) (interface{}, error) {
        return formatAnalysisResult(input), nil
    }))
    
    // 连接节点
    workflow.Connect(parseNode, extractNode)
    workflow.Connect(extractNode, analysisNode)
    workflow.Connect(analysisNode, formatNode)
    
    return workflow, nil
}
```

## 注意事项

1. **资源管理**: 确保工作流执行完成后正确释放资源
2. **并发安全**: 如果工作流会被并发执行，确保组件是并发安全的
3. **超时控制**: 为长时间运行的工作流设置合理的超时时间
4. **错误恢复**: 实现适当的错误恢复和重试机制
5. **性能优化**: 对于复杂工作流，考虑使用缓存和并行处理

## 相关资源

- [Eino 官方文档](https://github.com/cloudwego/eino)
- [Eino 示例代码](https://github.com/cloudwego/eino/tree/main/examples)
- [组件文档](https://github.com/cloudwego/eino-ext)

