package create

import (
	"github.com/urfave/cli"
	"github.com/maleck13/templator/model"
	"github.com/maleck13/templator/service"
)

func CreateTemplateCmd() cli.Command {
	return cli.Command{
		Name:      "app_template",
		ArgsUsage: "<name> --target=[openshift,kubernetes]",
		Usage:     "<name> --target=openshift",
		Action: func(context *cli.Context) error {
			if len(context.Args()) != 1 {
				return cli.NewExitError("expected one arg "+context.Command.Usage, 1)
			}
			if err := CreateTemplateAction(context.Args()[0], flag_Target); err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			return nil
		},
	}
}

func CreateTemplateAction(name, target string) error {
	templateService := service.NewTemplateService("local")
	template := model.NewApplicationTemplate(name)
	return templateService.SaveTemplate(name, template)
}
