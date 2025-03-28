package graph

import (
	"github.com/tmc/langchaingo/llms"
	"golang.org/x/net/context"
)

type GraphCallback interface {
	HandleNodeStart(ctx context.Context, node string, initialState []llms.MessageContent)
	HandleNodeEnd(ctx context.Context, node string, finalState []llms.MessageContent)
	HandleNodeStream(ctx context.Context, node string, chunk []byte)
	HandleEdgeEntry(ctx context.Context, edge string, initialState []llms.MessageContent)
	HandleEdgeExit(ctx context.Context, edge string, finalState []llms.MessageContent, output string)
}

type SimpleCallback struct{}

func (callback SimpleCallback) HandleNodeStart(ctx context.Context, node string, initialState []llms.MessageContent) {
}
func (callback SimpleCallback) HandleNodeEnd(ctx context.Context, node string, finalState []llms.MessageContent) {
}
func (callback SimpleCallback) HandleNodeStream(ctx context.Context, node string, chunk []byte) {
}
func (callback SimpleCallback) HandleEdgeEntry(ctx context.Context, edge string, initialState []llms.MessageContent) {
}
func (callback SimpleCallback) HandleEdgeExit(ctx context.Context, edge string, finalState []llms.MessageContent, output string) {
}
