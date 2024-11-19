package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools/duckduckgo"
	"github.com/tmc/langgraphgo/graph"
)

func main() {

	model, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		panic(err)
	}

	intialState := graph.NewMessagesState([]llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are an agent that has access to a Duck Duck go search engine. Please provide the user with the information they are looking for by using the search tool provided."),
	})

	tools := []llms.Tool{
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "search",
				Description: "Preforms Duck Duck Go web search",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "The search query",
						},
					},
				},
			},
		},
	}

	agent := func(ctx context.Context, state graph.MessagesState) (graph.MessagesState, error) {
		response, err := model.GenerateContent(ctx, state.Messages, llms.WithTools(tools))
		if err != nil {
			return state, err
		}

		if len(response.Choices[0].ToolCalls) > 0 {
			// log.Printf("Generating %v tool calls", len(response.Choices[0].ToolCalls))
			state.ToolCalls = response.Choices[0].ToolCalls

			msg := llms.TextParts(llms.ChatMessageTypeAI, response.Choices[0].Content)
			for _, toolCall := range state.ToolCalls {
				msg.Parts = append(msg.Parts, toolCall)
			}

			state.Messages = append(state.Messages, msg)

			// log.Printf("agent state: %s", state.Messages)
			return state, nil
		}

		msg := llms.TextParts(llms.ChatMessageTypeAI, response.Choices[0].Content)

		state.Messages = append(state.Messages, msg)
		return state, nil
	}

	search := func(ctx context.Context, state graph.MessagesState) (graph.MessagesState, error) {
		for index, toolCall := range state.ToolCalls {
			if toolCall.FunctionCall.Name == "search" {

				var args struct {
					Query string `json:"query"`
				}

				if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
					return state, err
				}

				search, err := duckduckgo.New(1, duckduckgo.DefaultUserAgent)
				if err != nil {
					log.Printf("search error: %v", err)
					return state, err
				}

				toolResponse, err := search.Call(ctx, args.Query)
				if err != nil {
					log.Printf("search error: %v", err)
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

				state.Messages = append(state.Messages, msg)
				state.ToolCalls = append(state.ToolCalls[:index], state.ToolCalls[index+1:]...)
			}
		}

		// log.Printf("search state messages: %v", state.Messages)
		return state, nil
	}

	useSearch := func(ctx context.Context, state graph.MessagesState) string {
		for _, toolCall := range state.ToolCalls {
			if toolCall.FunctionCall.Name == "search" {
				return "search"
			}
		}

		return graph.END
	}

	workflow := graph.NewMessageGraph()

	workflow.AddNode("agent", agent)
	workflow.AddNode("search", search)

	workflow.SetEntryPoint("agent")
	workflow.AddConditionalEdge("agent", useSearch)
	workflow.AddEdge("search", "agent")

	app, err := workflow.Compile()
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	intialState.Messages = append(
		intialState.Messages,
		llms.TextParts(llms.ChatMessageTypeHuman, "Who is the founder of Apple?"),
	)

	response, err := app.Invoke(context.Background(), intialState)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	lastMsg := response.Messages[len(response.Messages)-1]
	log.Printf("last msg: %v", lastMsg.Parts[0])
}
