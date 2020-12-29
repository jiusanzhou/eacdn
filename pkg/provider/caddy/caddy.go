package caddy

import (
	"go.zoe.im/eacdn/pkg/core"
	"go.zoe.im/eacdn/pkg/provider"
)

type caddy struct {
	provider.Options
}

func (s caddy) Init(opts ...provider.Option) error {

	return nil
}

func (s caddy) Create(site *core.Site, opts ...provider.Option) error {

	return nil
}

func (s caddy) Update(site *core.Site, opts ...provider.Option) error {

	return nil
}

func init() {
	provider.Register(caddy{}, "caddy")
}
