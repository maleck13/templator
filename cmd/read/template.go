package read

import (
	"encoding/json"
	"os"
	"text/template"

	"github.com/urfave/cli"
	"github.com/maleck13/templator/service"
)

const LIST_TEMPLATES_TEMPLATE = `
{{range $k,$v := .}}
 | {{$k}} |
{{end}}
`

func ReadTemplateCmd() cli.Command {
	return cli.Command{
		Name:      "app_template",
		ArgsUsage: "[name]",
		Usage:     "[name]",
		Action: func(context *cli.Context) error {
			if len(context.Args()) != 1 {
				return ListTemplateAction()
			}
			if err := ReadTemplateAction(context.Args()[0]); err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			return nil
		},
	}
}

func ListTemplateAction() error {
	templateService := service.NewTemplateService("local")
	data, err := templateService.ListTemplates()
	if err != nil {
		return cli.NewExitError("failed to load templates "+err.Error(), 1)
	}
	t := template.New("templatesList")
	outT, err := t.Parse(LIST_TEMPLATES_TEMPLATE)
	if err != nil {
		return cli.NewExitError("failed to parse template "+err.Error(), 1)
	}
	if err := outT.Execute(os.Stdout, data); err != nil {
		return cli.NewExitError("failed to execute template "+err.Error(), 1)
	}
	return nil
}

func ReadTemplateAction(name string) error {
	templateService := service.NewTemplateService("local")
	appTemp, err := templateService.GetTemplate(name)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	decoder := json.NewDecoder(os.Stdout)

	if err := decoder.Decode(appTemp); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}
