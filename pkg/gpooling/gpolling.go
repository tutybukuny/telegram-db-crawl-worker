package gpooling

import (
	"github.com/panjf2000/ants/v2"

	"crawl-worker/pkg/l"
)

// Pool - pooling struct
type Pool struct {
	antsPool *ants.Pool
}

// IPool - pooling interface
type IPool interface {
	Submit(task func())
	Release()
	Running() int
}

// New - init pooling
func New(maxPoolSize int, logger l.Logger) *Pool {
	pool, err := ants.NewPool(maxPoolSize, ants.WithNonblocking(false), ants.WithPanicHandler(func(data interface{}) {
		logger.Error("gpool execution error", l.Error(data.(error)))
	}))
	if err != nil {
		logger.Fatal("error when init gpooling", l.Error(err))
	}
	return &Pool{
		antsPool: pool,
	}
}

// Release - release all gorotine
func (p *Pool) Release() {
	p.antsPool.Release()
}

// Running - returns the number of the currently running goroutines.
func (p *Pool) Running() int {
	return p.antsPool.Running()
}

// Submit - submit a task to this pool
func (p *Pool) Submit(task func()) {
	p.antsPool.Submit(task)
}
