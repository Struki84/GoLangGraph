package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Struki84/GoLangGraph/graph"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {

	model, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		panic(err)
	}

	initialState := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a helpful assistant"),
	}

	agent := func(ctx context.Context, state []llms.MessageContent, opts ...graph.GraphOptions) ([]llms.MessageContent, error) {
		options := graph.Options{}
		for _, opt := range opts {
			opt(&options)
		}

		response, err := model.GenerateContent(ctx, state, llms.WithStreamingFunc(options.StreamHandler))
		if err != nil {
			return nil, err
		}

		state = append(state, llms.TextParts(llms.ChatMessageTypeAI, response.Choices[0].Content))
		return state, nil
	}

	workflow := graph.NewMessageGraph()

	workflow.AddNode("agent", agent)
	workflow.AddEdge("agent", graph.END)
	workflow.SetEntryPoint("agent")

	app, err := workflow.Compile()
	if err != nil {
		log.Println(err)
		return
	}

	initialState = append(
		initialState,
		llms.TextParts(llms.ChatMessageTypeAI, "Hello! How are you doing?"),
	)

	streamFunc := func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}

	_, err = app.Invoke(context.Background(), initialState, graph.WithStreamHandler(streamFunc))
	if err != nil {
		log.Println(err)
		return
	}
}
