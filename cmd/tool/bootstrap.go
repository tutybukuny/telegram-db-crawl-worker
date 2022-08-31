package main

import (
	"context"

	"crawl-worker/config"
	toolservice "crawl-worker/internal/service/tool"
	"crawl-worker/pkg/container"
	handleossignal "crawl-worker/pkg/handle-os-signal"
	"crawl-worker/pkg/l"
	"crawl-worker/pkg/telegram"
	"github.com/zelenin/go-tdlib/client"
)

func bootstrap(cfg *config.Config) {
	var ll l.Logger
	container.NamedResolve(&ll, "ll")
	var shutdown handleossignal.IShutdownHandler
	container.NamedResolve(&shutdown, "shutdown")

	_, cancel := context.WithCancel(context.Background())
	shutdown.HandleDefer(cancel)

	//region config
	//endregion

	//region init store
	//endregion

	//region init agent
	tdClient := telegram.New(cfg.TelegramConfig)
	container.NamedSingleton("tdClient", func() *client.Client {
		return tdClient
	})
	//endregion

	//region init service
	container.NamedSingleton("toolService", func() toolservice.IService {
		return toolservice.New()
	})
	//endregion
}
