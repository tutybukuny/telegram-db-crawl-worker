package l

import (
	"testing"

	"crawl-worker/pkg/l/sentry"
)

func TestNew(t *testing.T) {

	//ll = New()
	ll = NewWithSentry(&sentry.Configuration{
		DSN:   "https://6c823523782944c597fcc102c8b6ae4e@o390151.ingest.sentry.io/5231166",
		Trace: struct{ Disabled bool }{Disabled: false},
	})
	defer ll.Sync()
	a := map[string]interface{}{
		"testdebug": 1,
	}
	ll.Debug("test debug", Any("test debug", a))
	ll.Info("test info", Any("test debug", a))
	ll.Warn("test warn")
	//ll.Panic("fatal")
	ll.Error("test err")

}
