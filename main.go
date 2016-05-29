package main

import (
	"encoding/json"
	"os"

	"fmt"
	"io"

	"github.com/urfave/cli"
	"github.com/maleck13/templator/cmd/create"
	"github.com/maleck13/templator/cmd/del"
	"github.com/maleck13/templator/cmd/read"
	"github.com/maleck13/templator/model"
	"k8s.io/kubernetes/pkg/runtime"
)

var (
	nodes        int
	storage      bool
	nodeSelector bool
)

const (
	TEMPLATE_DB = "./.templates.json"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		create.CreateCmd(),
		del.DeleteCmd(),
		read.ReadCmd(),
		generateCmd(),
	}

	app.Run(os.Args)
}

func generateCmd() cli.Command {
	return cli.Command{
		Name:      "generate",
		ArgsUsage: "<template>",
		Action:    generateAction,
		Usage:     "generate <template> --nodes=3 --storage --nodeSelector",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:        "nodes",
				Destination: &nodes,
			},
			cli.BoolFlag{
				Name:        "storage",
				Destination: &storage,
			},
			cli.BoolFlag{
				Name:        "nodeSelector",
				Destination: &nodeSelector,
			},
		},
	}
}

type writeJson struct {
	Encoder *json.Encoder
	writer  io.Writer
}

func NewWriteJson() *writeJson {
	return &writeJson{
		json.NewEncoder(os.Stdout),
		os.Stdout,
	}
}

func (wj *writeJson) Write(p []byte) (n int, err error) {
	//todo find a way to format
	return wj.writer.Write(p)
}

func generateAction(context *cli.Context) error {
	if len(context.Args()) != 1 {
		return cli.NewExitError(context.Command.Usage, 1)
	}
	var templateName = context.Args()[0]
	file, err := os.Open(TEMPLATE_DB)
	if err != nil {
		return cli.NewExitError("failed to open file "+err.Error(), 1)
	}
	defer file.Close()
	var modelTemp map[string]*model.ApplicationTemplate

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&modelTemp); err != nil {
		return cli.NewExitError("failed to decode json "+err.Error(), 1)
	}

	appTemplate := modelTemp[templateName]
	osTemplate := &model.Template{}
	osTemplate.Kind = appTemplate.Kind
	osTemplate.APIVersion = appTemplate.APIVersion
	osTemplate.ObjectMeta = appTemplate.ObjectMeta

	for _, v := range appTemplate.DeploymentConfigs {
		preparedConfig := buildDeploymentConfigs(v)
		osTemplate.Objects = append(osTemplate.Objects, preparedConfig...)
	}

	for _, v := range appTemplate.Services {
		osTemplate.Objects = append(osTemplate.Objects, v)
	}

	data, err := json.MarshalIndent(osTemplate, "", " ")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	fmt.Println(string(data))

	return nil
}

//func buildTemplate(ostTemplate *ApplicationTemplate)*Template{
//
//}

func buildDeploymentConfigs(dc *model.OSTDeploymentConfig) []runtime.Object {
	builtConfigs := make([]runtime.Object, 0)
	if !storage {
		//remove volumes
		dc.Spec.Template.Spec.Volumes = nil
		for i := 0; i < len(dc.Spec.Template.Spec.Containers); i++ {
			dc.Spec.Template.Spec.Containers[i].VolumeMounts = nil
		}
	}
	if !nodeSelector {
		//remove nodeSelector
		dc.Spec.Template.Spec.NodeSelector = nil
	}
	if dc.Spec.DeploymentStrategy == model.DeploymentStrategy_PerNodeConfig {
		for i := 0; i < nodes; i++ {
			//clone the object  it is only a shallow clone
			var cloneDC *model.OSTDeploymentConfig = &model.OSTDeploymentConfig{}
			*cloneDC = *dc
			//need to replace some fields manually added for now but walking the object might be an option
			cloneDC.ObjectMeta.Name = fmt.Sprintf(cloneDC.ObjectMeta.Name, i)
			if storage {
				//inc name of claims
				for k := 0; k < len(dc.Spec.Template.Spec.Volumes); k++ {
					claimName := cloneDC.Spec.Template.Spec.Volumes[k].PersistentVolumeClaim.ClaimName
					cloneDC.Spec.Template.Spec.Volumes[k].PersistentVolumeClaim.ClaimName = fmt.Sprintf(claimName, k)
				}
			}
			if dc.Spec.ReplicaStrategy == model.ReplicationStrategy_EqualToNodes {
				cloneDC.Spec.Replicas = nodes

			}
			builtConfigs = append(builtConfigs, cloneDC)

		}
	} else if dc.Spec.DeploymentStrategy == model.DeploymentStrategy_SingleConfig {
		if dc.Spec.ReplicaStrategy == model.ReplicationStrategy_EqualToNodes {
			dc.Spec.Replicas = nodes
		}
		builtConfigs = append(builtConfigs, dc)
	} else {
		builtConfigs = append(builtConfigs, dc)
	}
	return builtConfigs
}
