package main

import (
	"context"

	"github.com/zelenin/go-tdlib/client"
	"gorm.io/gorm"

	"crawl-worker/config"
	"crawl-worker/internal/model/entity"
	"crawl-worker/internal/pkg/config/dbsaverconfig"
	channelrepo "crawl-worker/internal/repository/channel"
	"crawl-worker/internal/repository/media-message"
	"crawl-worker/internal/service/dbsaverservice"
	channelstore "crawl-worker/internal/storage/mysql/channel"
	mediamessagestore "crawl-worker/internal/storage/mysql/media-message"
	"crawl-worker/pkg/container"
	"crawl-worker/pkg/gpooling"
	handleossignal "crawl-worker/pkg/handle-os-signal"
	"crawl-worker/pkg/l"
	"crawl-worker/pkg/mysql"
	"crawl-worker/pkg/telegram"
	validator "crawl-worker/pkg/validator"
)

func bootstrap(cfg *config.Config) {
	var ll l.Logger
	container.NamedResolve(&ll, "ll")
	var shutdown handleossignal.IShutdownHandler
	container.NamedResolve(&shutdown, "shutdown")

	_, cancel := context.WithCancel(context.Background())
	shutdown.HandleDefer(cancel)

	container.NamedSingleton("gpooling", func() gpooling.IPool {
		return gpooling.New(cfg.MaxPoolSize, ll)
	})

	container.NamedSingleton("validator", func() validator.IValidator {
		return validator.New()
	})

	//region config
	rootConfig := dbsaverconfig.LoadConfig(cfg.ConfigFile)
	container.NamedSingleton("dbsaverConfigMap", func() dbsaverconfig.ConfigMap {
		return rootConfig.ConfigMap
	})
	container.NamedSingleton("filteredContents", func() []string {
		return rootConfig.FilteredContents
	})
	//endregion

	//region init store
	db := mysql.New(cfg.MysqlConfig, ll)
	mysql.AutoMigration(db, []any{
		&entity.MediaMessage{},
		&entity.Channel{},
	}, ll)

	container.NamedSingleton("db", func() *gorm.DB {
		return db
	})

	container.NamedSingleton("channelRepo", func() channelrepo.IRepo {
		return channelstore.New(db)
	})

	container.NamedSingleton("mediaMessageRepo", func() mediamessagerepo.IRepo {
		return mediamessagestore.New(db)
	})
	//endregion

	//region init agent
	tdClient := telegram.New(cfg.TelegramConfig)
	container.NamedSingleton("tdClient", func() *client.Client {
		return tdClient
	})
	//endregion

	//region init service
	container.NamedSingleton("dbsaverService", func() dbsaverservice.IService {
		return dbsaverservice.New()
	})
	//endregion
}
