package cli

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

// getGlobalFlags returns the global flags
func getGlobalFlags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:    "squeaky-lobster",
			Aliases: []string{"sl"},
			Value:   false,
			Hidden:  true,
		}),
	}
}
