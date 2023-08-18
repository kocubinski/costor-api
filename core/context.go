package core

import "context"

type Context struct {
	context.Context
	DryRun  bool
	HomeDir string
}
