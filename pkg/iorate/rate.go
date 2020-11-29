package iorate

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/jonboulle/clockwork"
)

type Rate struct {
	clock clockwork.Clock
	total uint64
	mu    sync.Mutex

	since time.Time
	prev  uint64
}

func New() *Rate {
	clock := clockwork.NewRealClock()
	return &Rate{
		clock: clock,
		since: clock.Now(),
		mu:    sync.Mutex{},
	}
}

func (r *Rate) Add(n int) {
	atomic.AddUint64(&r.total, uint64(n))
}

func (r *Rate) Calculate(per time.Duration) float64 {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.clock.Now()
	dur := now.Sub(r.since)

	delta := r.total - r.prev
	rate := float64(delta*uint64(per)) / float64(dur)

	r.prev = r.total
	r.since = now
	return rate
}
