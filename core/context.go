package core

import "context"

type Context struct {
	context.Context
	DryRun  bool
	HomeDir string
}

func FromContext(ctx context.Context) Context {
	c, ok := ctx.(Context)
	if !ok {
		panic("context is not core.Context")
	}
	return c
}
