// YOU CAN EDIT YOUR CUSTOM CONFIG HERE

package config

import (
	"crawl-worker/pkg/mysql"
	"crawl-worker/pkg/telegram"
)

// Config ...
//easyjson:json
type Config struct {
	Base         `mapstructure:",squash"`
	SentryConfig SentryConfig `json:"sentry" mapstructure:"sentry"`

	MysqlConfig    mysql.Config    `json:"mysql" mapstructure:"mysql"`
	TelegramConfig telegram.Config `json:"telegram" mapstructure:"telegram"`

	ConfigFile  string `json:"config_file" mapstructure:"config_file"`
	MaxPoolSize int    `json:"max_pool_size" mapstructure:"max_pool_size"`
	IsSaveRaw   bool   `json:"is_save_raw" mapstructure:"is_save_raw"`
}

// SentryConfig ...
type SentryConfig struct {
	Enabled bool   `json:"enabled" mapstructure:"enabled"`
	DNS     string `json:"dns" mapstructure:"dns"`
	Trace   bool   `json:"trace" mapstructure:"trace"`
}
