package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	einoTool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"

	"github.com/mozhou-tech/jssandbox-go/pkg/eino/tool"
)

func main() {
	// 设置日志
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	ctx := context.Background()

	// 1. 初始化 JSSandbox 工具
	// JSSandboxTool 实现了 eino 的 tool.InvokableTool 接口
	jsTool, err := tool.NewJSSandboxTool(ctx, &tool.JSSandboxConfig{
		DefaultTimeout: 0, // 使用默认 30s
	})
	if err != nil {
		logrus.Fatalf("创建 JSSandbox 工具失败: %v", err)
	}
	defer jsTool.Close()

	// 2. 初始化 LLM (使用 OpenAI 兼容接口)
	// 需要设置环境变量: OPENAI_API_KEY 和 OPENAI_BASE_URL
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		logrus.Warn("未设置 OPENAI_API_KEY 环境变量，程序可能无法正常运行")
	}
	baseURL := os.Getenv("OPENAI_BASE_URL")
	modelName := os.Getenv("OPENAI_MODEL")
	if modelName == "" {
		modelName = "gpt-4o"
	}

	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   modelName,
	})
	if err != nil {
		logrus.Fatalf("初始化 LLM 失败: %v", err)
	}

	// 获取工具信息并绑定到模型
	toolInfo, err := jsTool.Info(ctx)
	if err != nil {
		logrus.Fatalf("获取工具信息失败: %v", err)
	}
	chatModel.BindTools([]*schema.ToolInfo{toolInfo})

	// 3. 构建 Eino Graph
	// 为了维护对话历史，我们让图中流转的数据类型始终为 []*schema.Message
	graph := compose.NewGraph[[]*schema.Message, []*schema.Message]()

	// 首先创建工具执行节点组件
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []einoTool.BaseTool{jsTool},
	})
	if err != nil {
		logrus.Fatalf("创建工具节点失败: %v", err)
	}

	// 添加模型节点
	// 我们使用 Lambda 包装模型，以便在生成响应后将其追加到历史记录中
	chatLambda, err := compose.AnyLambda(
		func(ctx context.Context, input []*schema.Message, opts ...any) ([]*schema.Message, error) {
			res, err := chatModel.Generate(ctx, input)
			if err != nil {
				return nil, err
			}
			return append(input, res), nil
		},
		nil, nil, nil)
	if err != nil {
		logrus.Fatalf("创建模型 Lambda 失败: %v", err)
	}
	err = graph.AddLambdaNode("chat", chatLambda)
	if err != nil {
		logrus.Fatalf("添加模型节点失败: %v", err)
	}

	// 添加工具执行节点
	// 同样使用 Lambda 包装，执行最后一条消息中的工具调用，并将结果追加到历史记录中
	toolsLambda, err := compose.AnyLambda(
		func(ctx context.Context, input []*schema.Message, opts ...any) ([]*schema.Message, error) {
			if len(input) == 0 {
				return input, nil
			}
			lastMsg := input[len(input)-1]

			// toolsNode 接收 Assistant 消息并返回 Tool 消息列表
			toolResults, err := toolsNode.Invoke(ctx, lastMsg)
			if err != nil {
				return nil, err
			}

			return append(input, toolResults...), nil
		},
		nil, nil, nil)
	if err != nil {
		logrus.Fatalf("创建工具 Lambda 失败: %v", err)
	}
	err = graph.AddLambdaNode("tools", toolsLambda)
	if err != nil {
		logrus.Fatalf("添加工具节点失败: %v", err)
	}

	// 设置起点
	err = graph.AddEdge(compose.START, "chat")
	if err != nil {
		logrus.Fatalf("添加起始边失败: %v", err)
	}

	// 分支逻辑：检查最后一条消息是否有工具调用
	err = graph.AddBranch("chat", compose.NewGraphBranch(func(ctx context.Context, input []*schema.Message) (string, error) {
		if len(input) == 0 {
			return compose.END, nil
		}
		lastMsg := input[len(input)-1]
		if len(lastMsg.ToolCalls) > 0 {
			return "tools", nil
		}
		return compose.END, nil
	}, map[string]bool{"tools": true, compose.END: true}))
	if err != nil {
		logrus.Fatalf("添加分支失败: %v", err)
	}

	// 工具执行完后返回给模型
	err = graph.AddEdge("tools", "chat")
	if err != nil {
		logrus.Fatalf("添加工具回调边失败: %v", err)
	}

	// 编译图
	runnable, err := graph.Compile(ctx)
	if err != nil {
		logrus.Fatalf("编译图失败: %v", err)
	}

	// 4. 执行
	input := []*schema.Message{
		schema.SystemMessage("你是一个强大的助手，可以使用 jssandbox 工具执行 JavaScript 代码来完成任务。"),
		schema.UserMessage("请通过执行 JavaScript 代码获取当前的日期和时间。"),
	}

	fmt.Println("正在发送请求给 Agent...")
	output, err := runnable.Invoke(ctx, input)
	if err != nil {
		logrus.Fatalf("执行失败: %v", err)
	}

	fmt.Println("\n" + os.ExpandEnv("==================== Agent 输出 ===================="))
	// 输出最后一条消息的内容
	if len(output) > 0 {
		fmt.Printf("结果: %s\n", output[len(output)-1].Content)
	}
	fmt.Println(os.ExpandEnv("===================================================="))
}
