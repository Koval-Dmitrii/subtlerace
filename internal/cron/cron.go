package cron

import (
	"context"
	"time"
)

type Cron interface {
	Run(ctx context.Context, action func(), next func() time.Duration)
}

var _ Cron = (*cronImpl)(nil)

type cronImpl struct {
}

func New() *cronImpl {
	return &cronImpl{}
}

func (c *cronImpl) Run(ctx context.Context, action func(), next func() time.Duration) {
	d := next()

	t := time.NewTimer(d)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-t.C:
			select {
			case <-ctx.Done():
				return
			default:
			}

			action()

			nd := next()
			if nd < 0 {
				nd = 0
			}

			t.Reset(nd)
		}
	}
}
