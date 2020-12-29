package provider

import "go.zoe.im/eacdn/pkg/core"

type internal struct {
	Options
}

func (s internal) Init(opts ...Option) error {

	return nil
}

func (s internal) Create(site *core.Site, opts ...Option) error {

	return nil
}

func (s internal) Update(site *core.Site, opts ...Option) error {

	return nil
}

func init() {
	// TODO: implement self server with caddy
	// Register(internal{}, "internal")
}
