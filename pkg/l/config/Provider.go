package config

// ChangeCallback is called for updates of configuration data
type ChangeCallback func(key string, provider string, data interface{})

// Root marks the root node in a Provider
const Root = ""

// A Provider provides a unified interface to accessing
// configuration systems.
type Provider interface {
	// the Name of the provider (YAML, Env, etc)
	Name() string
	// Get pulls a config value
	Get(key string) Value

	// A RegisterChangeCallback provides callback registration for config providers.
	// These callbacks are nop if a dynamic provider is not configured for the service.
	RegisterChangeCallback(key string, callback ChangeCallback) error
	UnregisterChangeCallback(token string) error
}

// scopedProvider defines recursive interface of providers based on the prefix
type scopedProvider struct {
	Provider

	prefix string
}

// NewScopedProvider creates a child provider given a prefix
func NewScopedProvider(prefix string, provider Provider) Provider {
	if prefix == "" {
		return provider
	}

	return &scopedProvider{
		Provider: provider,
		prefix:   prefix,
	}
}

func (sp scopedProvider) addPrefix(key string) string {
	if key == "" {
		return sp.prefix
	}

	return sp.prefix + "." + key
}

// Get returns configuration value
func (sp scopedProvider) Get(key string) Value {
	return sp.Provider.Get(sp.addPrefix(key))
}

// RegisterChangeCallback registers the callback in the underlying provider
func (sp scopedProvider) RegisterChangeCallback(key string, callback ChangeCallback) error {
	return sp.Provider.RegisterChangeCallback(sp.addPrefix(key), callback)
}

// UnregisterChangeCallback un registers a callback in the underlying provider
func (sp scopedProvider) UnregisterChangeCallback(key string) error {
	return sp.Provider.UnregisterChangeCallback(sp.addPrefix(key))
}
