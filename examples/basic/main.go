package main

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langgraphgo/graph"
)

func main() {

	model, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		panic(err)
	}

	oracleNode := func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		reponse, err := model.GenerateContent(ctx, state, llms.WithTemperature(0.7))
		if err != nil {
			return nil, err
		}

		msg := llms.TextParts(schema.ChatMessageTypeAI, reponse.Choices[0].Content)

		return append(state, msg), nil
	}

	workflow := graph.NewMessageGraph()

	workflow.AddNode("oracle", oracleNode)
	workflow.AddNode(graph.END, func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		return state, nil
	})

	workflow.SetEntryPoint("oracle")
	workflow.AddEdge("oracle", graph.END)

	runable, err := workflow.Compile()
	if err != nil {
		panic(err)
	}

	intialState := []llms.MessageContent{
		llms.TextParts(schema.ChatMessageTypeSystem, "You are a helpful assistant."),
		llms.TextParts(schema.ChatMessageTypeHuman, "Hello, how are you?"),
	}

	response, err := runable.Invoke(context.Background(), intialState)
	if err != nil {
		panic(err)
	}

	log.Printf("response: %s", response)
}
