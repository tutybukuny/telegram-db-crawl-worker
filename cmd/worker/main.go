package main

import (
	"crawl-worker/config"
	"crawl-worker/pkg/container"
	handleossignal "crawl-worker/pkg/handle-os-signal"
	"crawl-worker/pkg/l"
	"crawl-worker/pkg/l/sentry"
)

func main() {
	ll := l.New()
	cfg := config.Load(ll)

	if cfg.SentryConfig.Enabled {
		ll = l.NewWithSentry(&sentry.Configuration{
			DSN: cfg.SentryConfig.DNS,
			Trace: struct{ Disabled bool }{
				Disabled: !cfg.SentryConfig.Trace,
			},
		})
	}

	container.NamedSingleton("ll", func() l.Logger {
		return ll
	})

	// init os signal handle
	shutdown := handleossignal.New(ll)
	shutdown.HandleDefer(func() {
		ll.Sync()
	})
	container.NamedSingleton("shutdown", func() handleossignal.IShutdownHandler {
		return shutdown
	})

	bootstrap(cfg)

	go startWorker(cfg)

	// handle signal
	if cfg.Environment == "D" {
		shutdown.SetTimeout(1)
	}
	shutdown.Handle()
}
