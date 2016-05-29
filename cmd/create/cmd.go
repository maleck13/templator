package create

import "github.com/urfave/cli"

var (
	flag_Target string
)

func CreateCmd() cli.Command {
	return cli.Command{
		Name: "create",
		Subcommands: []cli.Command{
			CreateTemplateCmd(),
			CreateDeploymentCmd(),
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "target",
				Usage:       "--target=kubernetes sets the template type to generate",
				Destination: &flag_Target,
			},
		},
	}
}
