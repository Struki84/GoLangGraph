package graph

import (
	"context"

	"github.com/tmc/langchaingo/llms"
)

const (
	DeafaultMessagesKey = "messages"
)

// MessageGraph represents a message graph.
type MessageGraph struct {
	Graph
	MessagesKey string
}

// NewMessageGraph creates a new instance of MessageGraph.
func NewMessageGraph(key ...string) *MessageGraph {
	graph := &MessageGraph{
		MessagesKey: DeafaultMessagesKey,
		Graph: Graph{
			nodes: make(map[string]Node),
		},
	}

	if len(key) > 0 {
		graph.MessagesKey = key[0]
	}

	return graph
}

// AddNode adds a new node to the message graph with the given name and function.
func (g *MessageGraph) AddNode(name string, fn func(ctx context.Context, messages []llms.MessageContent) ([]llms.MessageContent, error)) {

	g.nodes[name] = Node{
		Name: name,
		Function: func(ctx context.Context, state map[string]any) (map[string]any, error) {
			msgs, err := fn(ctx, state["messages"].([]llms.MessageContent))
			state[g.MessagesKey] = msgs

			return state, err
		},
	}
}

// AddEdge adds a new edge to the message graph between the "from" and "to" nodes.
func (g *MessageGraph) AddEdge(from, to string) {
	g.edges = append(g.edges, Edge{
		From: from,
		To: func(ctx context.Context, state map[string]any) string {
			return to
		},
	})
}

func (g *MessageGraph) AddConditionalEdge(from string, condition func(ctx context.Context, state []llms.MessageContent) string) {
	g.edges = append(g.edges, Edge{
		From: from,
		To: func(ctx context.Context, state map[string]any) string {
			return condition(ctx, state[g.MessagesKey].([]llms.MessageContent))
		},
	})
}

// SetEntryPoint sets the entry point node name for the message graph.
func (g *MessageGraph) SetEntryPoint(name string) {
	g.entryPoint = name
}

// Compile compiles the message graph and returns a Runnable instance.
// It returns an error if the entry point is not set.
func (g *MessageGraph) Compile() (*Runnable, error) {
	if g.entryPoint == "" {
		return nil, ErrEntryPointNotSet
	}

	return &Runnable{
		graph: &g.Graph,
	}, nil
}
