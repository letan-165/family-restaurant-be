package utils

import (
	"context"
	"time"
)

const DefaultTimeout = 10 * time.Second

func DefaultCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}
