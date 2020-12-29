package config

import (
	"go.zoe.im/eacdn/pkg/core"
	"go.zoe.im/eacdn/pkg/provider"
)

// Config contains ...
type Config struct {
	Addr      string            `json:"addr,omitempty" yaml:"addr"`
	Debug     bool              `json:"debug,omitempty" yaml:"debug"`
	Token     string            `json:"token,omitempty" yaml:"token"`
	Providers []provider.Object `json:"providers,omitempty" yaml:"providers" opts:"-"`
	Sites     []core.Site       `json:"sites,omitempty" yaml:"sites" opts:"-"`
}

// New ...
func New() *Config {
	return &Config{
		Addr:  ":30911",
		Debug: true,
		Token: "eacdn-2020!",
	}
}
