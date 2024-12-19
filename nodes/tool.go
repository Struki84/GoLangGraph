package nodes

import (
	"context"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

func ToolNode(nodeTools []tools.Tool) func(context.Context, []llms.MessageContent) ([]llms.MessageContent, error) {
	return func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		lastMsg := state[len(state)-1]

		for _, part := range lastMsg.Parts {
			toolCall, ok := part.(llms.ToolCall)

			if ok {
				for _, tool := range nodeTools {
					if tool.Definition().Name == toolCall.FunctionCall.Name {
						toolResponse, err := tool.Call(ctx, toolCall.FunctionCall.Arguments)
						if err != nil {
							return state, err
						}

						msg := llms.MessageContent{
							Role: llms.ChatMessageTypeTool,
							Parts: []llms.ContentPart{
								llms.ToolCallResponse{
									ToolCallID: toolCall.ID,
									Name:       toolCall.FunctionCall.Name,
									Content:    toolResponse,
								},
							},
						}

						state = append(state, msg)
						return state, nil
					}
				}
			}
		}

		return state, nil
	}
}
