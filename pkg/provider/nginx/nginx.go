package nginx

import (
	"go.zoe.im/eacdn/pkg/core"
	"go.zoe.im/eacdn/pkg/provider"
)

type nginx struct {
	provider.Options
}

func (s nginx) Init(opts ...provider.Option) error {
	// make sure nginx command exits
	// make sure nginx process alive
	return nil
}

func (s nginx) Create(site *core.Site, opts ...provider.Option) error {

	return nil
}

func (s nginx) Update(site *core.Site, opts ...provider.Option) error {

	return nil
}

func init() {
	provider.Register(nginx{}, "nginx")
}
