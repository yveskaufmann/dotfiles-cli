package provider

import (
	"sort"

	"yv35.com/dotfiles/internal/types"
)

var defaultProviders = []types.Provider{}

func RegisterProviders(p types.Provider) {
	defaultProviders = append(defaultProviders, p)
}

type Registry struct {
	Providers map[string]types.Provider
}

func NewRegistry() *Registry {

	providers := make(map[string]types.Provider)
	for _, p := range defaultProviders {
		providers[p.ID()] = p
	}

	return &Registry{
		Providers: providers,
	}
}

func (r *Registry) RegisterProvider(provider types.Provider) {
	r.Providers[provider.ID()] = provider
}

func (r *Registry) GetProvider(id string) (types.Provider, bool) {
	provider, exists := r.Providers[id]
	return provider, exists
}

func (r *Registry) List() []types.Provider {
	providers := make([]types.Provider, 0, len(r.Providers))
	for _, provider := range r.Providers {
		providers = append(providers, provider)
	}

	// Sort by priority (lower numbers first)
	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Priority() < providers[j].Priority()
	})

	return providers
}
