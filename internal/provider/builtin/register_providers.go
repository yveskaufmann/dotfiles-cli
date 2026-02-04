package builtin

import (
	"yv35.com/dotfiles/internal/provider"
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
)

func registerProviders() {
	provider.RegisterProviders(github.NewProvider())
	provider.RegisterProviders(pipx.NewProvider())
	provider.RegisterProviders(apt.NewProvider())
	provider.RegisterProviders(brew.NewProvider())
	provider.RegisterProviders(nvm.NewProvider())
	provider.RegisterProviders(sdkman.NewProvider())
	provider.RegisterProviders(npm.NewProvider())
	provider.RegisterProviders(snap.NewProvider())
	provider.RegisterProviders(binary.NewProvider())
	provider.RegisterProviders(script.NewProvider())
	provider.RegisterProviders(custom.NewProvider())
	provider.RegisterProviders(example.NewProvider())
}

func init() {
	registerProviders()
}
