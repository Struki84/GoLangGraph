package graph

import "context"

type GraphOptions func(*Options)

type Options struct {
	StreamHandler   func(ctx context.Context, chunk []byte) error
	CallbackHandler GraphCallback
}

func WithStreamHandler(handler func(ctx context.Context, chunk []byte) error) GraphOptions {
	return func(opts *Options) {
		opts.StreamHandler = handler
	}
}

func WithCallback(callback GraphCallback) GraphOptions {
	return func(opts *Options) {
		opts.CallbackHandler = callback
	}
}
