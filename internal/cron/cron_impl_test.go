package cron

import "os"

func newCronForTest() Cron {
	if os.Getenv("CRON_IMPL") == "afterfunc" {
		return NewAfterFunc()
	}

	return New()
}
