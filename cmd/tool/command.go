package main

import (
	"crawl-worker/internal/commander"
	"crawl-worker/pkg/container"
	handleossignal "crawl-worker/pkg/handle-os-signal"
	"crawl-worker/pkg/l"
)

func runCmd() {
	var ll l.Logger
	container.NamedResolve(&ll, "ll")
	var shutdown handleossignal.IShutdownHandler
	container.NamedResolve(&shutdown, "shutdown")

	cmd := commander.New()
	cmd.Execute()
}
