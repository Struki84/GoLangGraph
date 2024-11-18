package main

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langgraphgo/graph"
)

func main() {

	model, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		panic(err)
	}

	intialState := make(map[string]any)
	intialState["messages"] = []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a helpful assistant."),
	}

	tools := []llms.Tool{
		{
			Type: "Function",
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

	agent := func(ctx context.Context, state map[string]any) (map[string]any, error) {
		msgs := state["messages"].([]llms.MessageContent)

		reponse, err := model.GenerateContent(ctx, msgs, llms.WithTools(tools))
		if err != nil {
			return nil, err
		}
		state["content"] = reponse

		msg := llms.TextParts(llms.ChatMessageTypeAI, reponse.Choices[0].Content)
		state["messages"] = append(msgs, msg)

		log.Printf("msgs: %v", msgs)
		return state, nil
	}

	search := func(ctx context.Context, state map[string]any) (map[string]any, error) {
		return state, nil
	}

	useSearch := func(ctx context.Context, state map[string]any) string {
		response := state["content"].(llms.ContentResponse)
		for _, tool := range response.Choices[0].ToolCalls {
			if tool.FunctionCall.Name == "search" {
				return "search"
			}
		}

		return graph.END
	}

	workflow := graph.New()

	workflow.AddNode("agent", agent)
	workflow.AddNode("search", search)

	workflow.SetEntryPoint("agent")
	workflow.AddConditionalEdge("agent", useSearch)
	workflow.AddEdge("relay", graph.END)

	runable, _ := workflow.Compile()

	intialState["messages"] = append(
		intialState["messages"].([]llms.MessageContent),
		llms.TextParts(llms.ChatMessageTypeHuman, "Hello, how are you?"),
	)

	response, err := runable.Invoke(context.Background(), intialState)
	if err != nil {
		panic(err)
	}

	log.Printf("response: %s", response)
}
