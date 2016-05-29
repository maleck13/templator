package del

import "github.com/urfave/cli"

func DeleteCmd() cli.Command {
	return cli.Command{
		Name: "delete",
		Subcommands: []cli.Command{
			DeleteTemplateCmd(),
		},
	}
}
