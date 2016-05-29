package del

import (
	"github.com/urfave/cli"
	"github.com/maleck13/templator/service"
)

func DeleteTemplateCmd() cli.Command {
	return cli.Command{
		Name:      "app_template",
		ArgsUsage: "<name>",
		Usage:     "<name>",
		Action: func(context *cli.Context) error {
			if len(context.Args()) != 1 {
				return cli.NewExitError("expected one arg "+context.Command.Usage, 1)
			}
			if err := DeleteTemplateAction(context.Args()[0]); err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			return nil
		},
	}
}

func DeleteTemplateAction(name string) error {
	templateService := service.NewTemplateService("local")
	return templateService.DeleteTemplate(name)
}
