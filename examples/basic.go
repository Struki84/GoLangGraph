package examples

import (
	"context"

	"github.com/tmc/langchaingo/graph"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

func main() {

	model, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		panic(err)
	}

	oracleNode := func(ctx context.Context, state []graph.MessageContent) ([]graph.MessageContent, error) {
		reponse, err := model.GenerateContent(ctx, state, llms.WithTemperature(0.7))
		if err != nil {
			return nil, err
		}

		msg := llms.TextParts(schema.ChatMessageTypeAI, reponse.Choices[0].Content)

		return append(state, msg), nil
	}

	graph := graph.NewMessageGraph()

	graph.addNode("oracle", oracleNode)
	graph.addNode(graph.END, func(ctx context.Context, state []graph.MessageContent) ([]graph.MessageContent, error) {
		return state, nil
	})

	graph.SetEntryPoint("oracle")
	graph.AddEdge("oracle", graph.END)

	workflow := graph.Compile()
	
	response, err := workflow.Invoke(context.Background(), []graph.MessageContent{})
	if err != nil {
		panic(err)
	}

	println(response)
}
