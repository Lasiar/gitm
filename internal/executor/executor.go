package executor

import "context"

type Command struct {
	Script string
}

type Executor interface {
	Run(ctx context.Context, cmd Command) error
}
