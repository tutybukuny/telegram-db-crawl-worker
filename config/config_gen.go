//AUTO-GENERATED: DO NOT EDIT

package config

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"crawl-worker/pkg/l"
)

// Base ...
type Base struct {
	Environment string `json:"environment" mapstructure:"environment"  validate:"required"`
	LogLevel    string `json:"log_level" mapstructure:"log_level"`
	LogColor    bool   `json:"log_color" mapstructure:"log_color"`
}

// Load ...
func Load(ll l.Logger, cPath ...string) *Config {
	var cfg = &Config{}
	v := viper.NewWithOptions(viper.KeyDelimiter("__"))

	customConfigPath := "."
	if len(cPath) > 0 {
		customConfigPath = cPath[0]
	}

	v.SetConfigType("env")
	v.SetConfigFile(".env")
	if len(cPath) > 0 {
		v.SetConfigName(".env")
	}
	v.AddConfigPath(customConfigPath)
	v.AddConfigPath(".")
	v.AddConfigPath("/app")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		ll.Fatal("Error reading config file", l.Error(err))
	}

	err := v.Unmarshal(&cfg)
	if err != nil {
		ll.Fatal("Failed to unmarshal config", l.Error(err))
	}

	ll.Debug("Config loaded", l.Object("config", cfg))

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			ll.S.Fatalf("Invalid config [%+v], tag [%+v], value [%+v]", err.StructNamespace(), err.Tag(), err.Value())
		}
	}

	return cfg
}
