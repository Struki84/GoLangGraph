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

	intialState := make(map[string]any)
	intialState["messages"] = []llms.MessageContent{
		llms.TextParts(schema.ChatMessageTypeSystem, "You are a helpful assistant."),
	}

	agent := func(ctx context.Context, state map[string]any) (map[string]any, error) {
		msgs := state["messages"].([]llms.MessageContent)

		reponse, err := model.GenerateContent(ctx, msgs, llms.WithTemperature(0.7))
		if err != nil {
			return nil, err
		}

		msg := llms.TextParts(schema.ChatMessageTypeAI, reponse.Choices[0].Content)
		state["messages"] = append(msgs, msg)

		log.Printf("msgs: %v", msgs)
		return state, nil
	}

	relay := func(ctx context.Context, state map[string]any) (map[string]any, error) {
		log.Printf("I'm a Relay")
		return state, nil
	}

	useRelay := func(ctx context.Context, state map[string]any) string {

		return "relay"
	}

	workflow := graph.New()

	workflow.AddNode("agent", agent)
	workflow.AddNode("relay", relay)

	workflow.SetEntryPoint("agent")
	workflow.AddConditionalEdge("agent", useRelay)
	workflow.AddEdge("relay", graph.END)

	runable, _ := workflow.Compile()

	intialState["messages"] = []llms.MessageContent{
		llms.TextParts(schema.ChatMessageTypeHuman, "Hello, how are you?"),
	}

	response, err := runable.Invoke(context.Background(), intialState)
	if err != nil {
		panic(err)
	}

	log.Printf("response: %s", response)
}
