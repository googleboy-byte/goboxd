package stats

import (
	"sync/atomic"
	"time"
)

type Stats struct {
	InFlight           atomic.Int64
	JobsTotal          atomic.Int64
	JobsFailedInternal atomic.Int64
	LastInternalErrAt  atomic.Pointer[time.Time]
}

func NewStats() *Stats {
	return &Stats{}
}
