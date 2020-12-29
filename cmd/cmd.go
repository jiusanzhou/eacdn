package cmd

import (
	"go.zoe.im/x/cli"
	"go.zoe.im/x/version"
)

var (
	// root command to contains all sub commands
	cmd = cli.New(
		// set name and description in run function
		cli.Name("eacdn"),
		cli.Short("EaCDN is an simple CDN manager."),
		version.NewOption(true),
		cli.Run(func(c *cli.Command, args ...string) {
			c.Help()
		}),
	)
)

// Register sub command
func Register(scs ...*cli.Command) {
	cmd.Register(scs...)
}

// Run call the global's command run
func Run(opts ...cli.Option) error {
	return cmd.Run(opts...)
}

// Option reload with options
func Option(opts ...cli.Option) {
	cmd.Option(opts...)
}
