package cron

import (
	"context"
	"time"
)

type afterFuncCronImpl struct{}

var _ Cron = (*afterFuncCronImpl)(nil)

func NewAfterFunc() *afterFuncCronImpl {
	return &afterFuncCronImpl{}
}

func (c *afterFuncCronImpl) Run(ctx context.Context, action func(), next func() time.Duration) {
	var t *time.Timer

	t = time.AfterFunc(next(), func() {
		select {
		case <-ctx.Done():
			return
		default:
			action()
			t.Reset(next())
		}
	})

	<-ctx.Done()
}
