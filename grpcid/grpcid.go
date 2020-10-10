// Package grpcid provides a gRPC middleware interceptor that adds a string ID
// to the request context.
package grpcid

import (
	"context"

	ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"

	"github.com/kenshaw/stringid"
)

const (
	// DefaultTag is the default tag name used.
	DefaultTag = "correlation.id"
)

// UnaryServerInterceptor creates a unary server interceptor for gRPC that adds
// a string ID to the request context.
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	// create config
	c := &config{tag: DefaultTag}

	// apply opts
	for _, o := range opts {
		o(c)
	}

	// ensure generator exists
	if c.g == nil {
		c.g = stringid.NewPushGenerator(nil)
	}

	return func(ctxt context.Context, req interface{}, info *grpc.UnaryServerInfo, h grpc.UnaryHandler) (interface{}, error) {
		id := c.g.Generate()
		ctxtags.Extract(ctxt).Set(c.tag, id)
		return h(stringid.WithID(ctxt, id), req)
	}
}

// Option is a grpcid interceptor option.
type Option func(*config)

// config holds configurtion information about a grpcid interceptor.
type config struct {
	g   stringid.Generator
	tag string
}

// Generator is a grpcid interceptor option to set the string ID generator to use.
func Generator(g stringid.Generator) Option {
	return func(c *config) {
		c.g = g
	}
}

// Tag is a grpcid interceptor option to set the tag name used.
func Tag(tag string) Option {
	return func(c *config) {
		c.tag = tag
	}
}
