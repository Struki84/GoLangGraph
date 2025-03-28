package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Struki84/GoLangGraph/graph"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type MyCallback struct {
	graph.SimpleCallback
}

func (callback MyCallback) HandleNodeStart(ctx context.Context, node string, initialState []llms.MessageContent) {
	log.Println("Callback from node:", node)
}

func (callback MyCallback) HandleNodeStream(ctx context.Context, node string, chunk []byte) {
	fmt.Print(string(chunk))
}

func main() {

	model, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		panic(err)
	}

	initialState := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a helpful assistant"),
	}

	agent := func(ctx context.Context, state []llms.MessageContent, opts graph.Options) ([]llms.MessageContent, error) {

		opts.CallbackHandler.HandleNodeStart(ctx, "agent", state)

		response, err := model.GenerateContent(ctx, state,
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				opts.CallbackHandler.HandleNodeStream(ctx, "agent", chunk)
				return nil
			}),
		)
		if err != nil {
			return state, err
		}

		state = append(state, llms.TextParts(llms.ChatMessageTypeAI, response.Choices[0].Content))

		opts.CallbackHandler.HandleNodeEnd(ctx, "agent", state)
		return state, nil
	}
	callback := MyCallback{}
	workflow := graph.NewMessageGraph(graph.WithCallback(callback))

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

	_, err = app.Invoke(context.Background(), initialState)
	if err != nil {
		log.Println(err)
		return
	}
}
