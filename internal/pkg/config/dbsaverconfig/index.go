package dbsaverconfig

import (
	"io/ioutil"
	"os"

	"crawl-worker/pkg/container"
	"crawl-worker/pkg/json"
	"crawl-worker/pkg/l"
)

type Config struct {
	ChannelID   int64 `json:"channel_id,omitempty"`
	ChannelType int   `json:"channel_type,omitempty"` // match with media message type
}

type ConfigMap map[string]*Config

type RootConfig struct {
	ConfigMap        ConfigMap `json:"config_map,omitempty"`
	FilteredContents []string  `json:"filtered_contents,omitempty"`
	DBChannelID      int64     `json:"db_channel_id"`
}

func LoadConfig(configPath string) *RootConfig {
	var ll l.Logger
	container.NamedResolve(&ll, "ll")

	if configPath == "" {
		configPath = "./config.json"
	}

	rootConfig := new(RootConfig)

	file, err := os.Open(configPath)
	if err != nil {
		ll.Fatal("cannot read dbsaver config", l.String("config_path", configPath), l.Error(err))
	}
	defer file.Close()
	configJson, err := ioutil.ReadAll(file)
	if err != nil {
		ll.Fatal("cannot read dbsaver config", l.String("config_path", configPath), l.Error(err))
	}
	err = json.Unmarshal(configJson, rootConfig)
	if err != nil {
		ll.Fatal("cannot parse channel config",
			l.String("channel_config_file", configPath),
			l.ByteString("config_json", configJson), l.Error(err))
	}
	ll.Info("loaded dbsaver config", l.Object("channel_config", rootConfig))

	return rootConfig
}
