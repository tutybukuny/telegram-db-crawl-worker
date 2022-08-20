package main

import (
	"crawl-worker/config"
	"crawl-worker/internal/listener"
	"crawl-worker/pkg/container"
	handleossignal "crawl-worker/pkg/handle-os-signal"
	"crawl-worker/pkg/l"
)

func startWorker(cfg *config.Config) {
	var ll l.Logger
	container.NamedResolve(&ll, "ll")
	var shutdown handleossignal.IShutdownHandler
	container.NamedResolve(&shutdown, "shutdown")

	worker := listener.New()
	worker.Listen()
}
