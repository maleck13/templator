package read

import "github.com/urfave/cli"

func ReadCmd() cli.Command {
	return cli.Command{
		Name: "read",
		Subcommands: []cli.Command{
			ReadTemplateCmd(),
		},
	}
}
