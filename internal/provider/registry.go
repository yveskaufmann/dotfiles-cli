package provider

import (
	"sort"

	"yv35.com/dotfiles/internal/provider/apt"
	"yv35.com/dotfiles/internal/provider/binary"
	"yv35.com/dotfiles/internal/provider/brew"
	"yv35.com/dotfiles/internal/provider/custom"
	"yv35.com/dotfiles/internal/provider/example"
	"yv35.com/dotfiles/internal/provider/github"
	"yv35.com/dotfiles/internal/provider/npm"
	"yv35.com/dotfiles/internal/provider/nvm"
	"yv35.com/dotfiles/internal/provider/pipx"
	"yv35.com/dotfiles/internal/provider/script"
	"yv35.com/dotfiles/internal/provider/sdkman"
	"yv35.com/dotfiles/internal/provider/snap"
	"yv35.com/dotfiles/internal/types"
)

type Registry struct {
	Providers map[string]types.Provider
}

func NewRegistry() *Registry {
	return &Registry{
		Providers: map[string]types.Provider{
			"example": example.NewProvider(),
			"snap":    snap.NewProvider(),
			"pipx":    pipx.NewProvider(),
			"npm":     npm.NewProvider(),
			"brew":    brew.NewProvider(),
			"binary":  binary.NewProvider(),
			"github":  github.NewProvider(),
			"apt":     apt.NewProvider(),
			"script":  script.NewProvider(),
			"custom":  custom.NewProvider(),
			"nvm":     nvm.NewProvider(),
			"sdkman":  sdkman.NewProvider(),
		},
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
