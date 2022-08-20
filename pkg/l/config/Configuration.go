package config

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"crawl-worker/pkg/l/sentry"
)

// Configuration defines the desired logging options.
type Configuration struct {
	zap.Config

	Sentry *sentry.Configuration `yaml:"sentry"`
}

// Configure initializes logging configuration struct from config provider
func (c *Configuration) Configure(cfg Value) error {

	// Because log.Configuration embeds zap, the PopulateStruct
	// does not work properly as it's unable to serialize fields directly
	// into the embedded struct, so inner struct has to be treated as a
	// separate object
	//
	// first, use the default zap configuration
	zapCfg := DefaultConfiguration().Config

	// override the embedded zap.Config stuct from config
	if err := cfg.PopulateStruct(&zapCfg); err != nil {
		return errors.New("unable to parse logging config")
	}

	// use the overriden zap config
	c.Config = zapCfg

	// override any remaining things fom config, i.e. Sentry
	if err := cfg.PopulateStruct(&c); err != nil {
		return errors.New("unable to parse logging config")
	}

	return nil
}

// Build constructs a *zap.Logger with the configured parameters.
func (c Configuration) Build(opts ...zap.Option) (*zap.Logger, error) {
	logger, err := c.Config.Build(opts...)
	if err != nil || c.Sentry == nil {
		// If there's an error or there's no Sentry config, we don't need to do
		// anything but delegate.
		return logger, err
	}
	sentryObj, err := c.Sentry.Build()
	if err != nil {
		return logger, err
	}
	return logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, sentryObj)
	})), nil
}

// DefaultConfiguration returns a fallback configuration for applications that
// don't explicitly configure logging.
func DefaultConfiguration() Configuration {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}

	return Configuration{
		Config: cfg,
	}
}
