package example

import (
	"fmt"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/types"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) Install(config config.DependencyGroup, onComplete types.OnTaskComplete) error {
	// Example installation logic goes here
	return nil
}

func (p *Provider) ID() string {
	return "example"
}

func (p *Provider) Setup() error {
	fmt.Println("Setting up ExampleProvider...")
	return nil
}

func (p *Provider) TearDown() error {
	fmt.Println("Tearing down ExampleProvider...")
	return nil
}
