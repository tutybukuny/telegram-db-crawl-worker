package handleossignal

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"crawl-worker/pkg/l"
)

// IShutdownHandler ...
type IShutdownHandler interface {
	SetTimeout(t time.Duration)
	SetLogger(ll l.Logger)
	HandleDefer(f func())
	Handle()
}

type handler struct {
	closes             []func()
	mu                 sync.Mutex
	delayForceShutdown time.Duration
	ll                 l.Logger
}

// New ... new handle os signal
func New(ll l.Logger) *handler {
	return &handler{
		delayForceShutdown: 15,
		ll:                 ll,
	}
}

// SetTimeout ...set delay time force shutdown in second. default is 15s
func (o *handler) SetTimeout(t time.Duration) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.delayForceShutdown = t
}

// SetLogger ...
func (o *handler) SetLogger(ll l.Logger) {
	o.ll = ll
}

// HandleDefer ...register clean-up actions when shutdown
func (o *handler) HandleDefer(f func()) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.closes = append(o.closes, f)
}

// Handle ...waiting for a signal and do clean actions
func (o *handler) Handle() {
	// handle signal
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	<-osSignal

	o.ll.Info("Shutdown started...")
	count := len(o.closes)
	if count == 0 {
		o.ll.Info("Bye ^^")
		return
	}
	// doneCloses
	doneCloses := make(chan struct{})
	for i := len(o.closes) - 1; i >= 0; i-- {
		f := o.closes[i]
		go o.closeObj(f, doneCloses)
	}
	timer := time.NewTimer(o.delayForceShutdown * time.Second)
	for {
		select {
		case <-timer.C:
			o.ll.Fatal("Force shutdown due to timeout!")
		case <-doneCloses:
			count--
			if count > 0 {
				continue
			}
			o.ll.Info("Bye ^^")
			os.Exit(0)
			return
		}
	}
}

func (o *handler) closeObj(closer func(), done chan struct{}) {
	closer()
	done <- struct{}{}
}
