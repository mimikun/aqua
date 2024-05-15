package cli

import (
	"github.com/urfave/cli/v3"
)

func (r *Runner) newVersionCommand() *cli.Command {
	return &cli.Command{
		Name:   "version",
		Usage:  "Show version",
		Action: r.versionAction,
	}
}

func (r *Runner) versionAction(c *cli.Command) error {
	cli.ShowVersion(c)
	return nil
}
