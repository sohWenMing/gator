package helper

import (
	"context"
	"time"
)

type TimeoutContextWithCancelFunc struct {
	Context    context.Context
	CancelFunc context.CancelFunc
}

func SpawnTimeOutContext(parent context.Context, duration time.Duration) TimeoutContextWithCancelFunc {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancelFunc := context.WithTimeout(parent, duration)
	return TimeoutContextWithCancelFunc{
		ctx, cancelFunc,
	}
}
