package create

import (
	"github.com/urfave/cli"
	"github.com/maleck13/templator/model"
	"github.com/maleck13/templator/service"
	"log"
	k8 "k8s.io/kubernetes/pkg/api/v1"
	k8resources "k8s.io/kubernetes/pkg/api/resource"
	"strconv"
	"strings"
	"github.com/maleck13/templator/cmd"
	"fmt"
	"k8s.io/kubernetes/pkg/util/intstr"
)

//is deployment a good name? it is a replication controller or deployment config

func CreateDeploymentCmd() cli.Command {
	return cli.Command{
		Name:      "deployment",
		ArgsUsage: "<name> <template>",
		Usage:     "deployment <name> <template>",
		Action: func(context *cli.Context) error {
			if len(context.Args()) != 2 {
				return cli.NewExitError("expected two args "+context.Command.ArgsUsage, 1)
			}
			return CreateDeploymentAction(context.Args()[0], context.Args()[1])
		},
	}
}

func addServices(deploymentModel *model.OSTDeploymentConfig) (*k8.Service, error) {
	//wrap in constructor
	serviceTemp := &k8.Service{}
	serviceTemp.APIVersion = "v1"
	serviceTemp.Kind = "Service"
	serviceTemp.Spec.Selector = make(map[string]string)
	serviceTemp.Spec.Selector["name"] = deploymentModel.Name
	serviceTemp.Spec.Ports = make([]k8.ServicePort, 0)
	cmd.QuestionAndAnswer("Do you want to expose any services for this deployment : ", func(answer string) {
		if "n" == answer {
			return
		}

		cmd.QuestionAndAnswer("name the service : ", func(name string) {
			serviceTemp.ObjectMeta.Name = name
		})

		cmd.QuestionAndAnswer("which ports do you want to expose (8080,3000) : ", func(answer string) {
			ports := strings.Split(answer, ",")
			for i, p := range ports {
				port := k8.ServicePort{}
				port.Name = fmt.Sprintf("%s-port-%d", serviceTemp.Name, i)
				port.Protocol = "TCP"
				pN, _ := strconv.ParseInt(p, 10, 32) //fix ignored error
				port.Port = int32(pN)
				cmd.QuestionAndAnswer("what is the target port for "+p+" :", func(answer string) {
					pN, _ := strconv.ParseInt(answer, 10, 32) //fix ignored error
					port.TargetPort = intstr.IntOrString{Type: intstr.Int, IntVal: int32(pN)}
				})
				serviceTemp.Spec.Ports = append(serviceTemp.Spec.Ports, port)
			}

		})

	})
	return serviceTemp, nil
}

func addContainers(deploymentModel *model.OSTDeploymentConfig) {
	cmd.QuestionAndAnswer("Do you want to add a container (y/n) : ", func(answer string) {
		if "no" == answer || "n" == answer {
			return
		}
	})

	container := k8.Container{}
	k8.SetDefaults_Container(&container)

	container.SecurityContext = &k8.SecurityContext{}

	cmd.QuestionAndAnswer("What is the name :", func(answer string) {
		container.Name = answer
	})
	cmd.QuestionAndAnswer("What image do you want to use :", func(answer string) {
		container.Image = answer
	})
	cmd.QuestionAndAnswer("What ports do you want to expose (8080,8443) :", func(answer string) {
		ports := strings.Split(answer, ",")
		for _, p := range ports {
			i, err := strconv.ParseInt(p, 10, 32)
			if err != nil {
				log.Fatal("could not parse int ", err)
			}
			container.Ports = append(container.Ports, k8.ContainerPort{
				ContainerPort: int32(i),
				Protocol:      "TCP", // default for now
			})
		}
	})
	cmd.QuestionAndAnswer("Do you need to set resource limits?:", func(answer string) {
		if "y" == strings.ToLower(answer) {
			cmd.QuestionAndAnswer("What's the max cpu shares : ", func(answer string) {
				container.Resources = k8.ResourceRequirements{
					Limits: k8.ResourceList{
						k8.ResourceName("cpu"): k8resources.MustParse(answer),
					},
				}
			})
			cmd.QuestionAndAnswer("What's the min cpu shares : ", func(answer string) {
				container.Resources.Requests = k8.ResourceList{
					k8.ResourceName("cpu"): k8resources.MustParse(answer),
				}

			})
			cmd.QuestionAndAnswer("What is the max memory resources : ", func(answer string) {
				container.Resources.Limits[k8.ResourceName("memory")] = k8resources.MustParse(answer)
			})
		}
	})
	cmd.QuestionAndAnswer("Any env vars? (MY_ENV_VAR:MY_VALUE,MY_ENV_TWO:MY_VAL_TWO)", func(answer string) {
		if "" == answer {
			return
		}
		envs := strings.Split(",", answer)
		fmt.Println(envs)

		for _, e := range envs {
			keyVal := strings.Split(",", e) //handle out of range prob
			container.Env = append(container.Env, k8.EnvVar{
				Name:  keyVal[0],
				Value: keyVal[1],
			})
		}
	})
	cmd.QuestionAndAnswer("Want to add another container ? (y/n) ", func(answer string) {
		if "y" == answer {
			addContainers(deploymentModel)
			return
		}
	})
	deploymentModel.Spec.Template.Spec.Containers = append(deploymentModel.Spec.Template.Spec.Containers, container)

}

func CreateDeploymentAction(name, temp string) error {
	var deploymentModel = model.NewOstDeploymentConfig(name)

	templateServ := service.NewTemplateService("local")

	addContainers(deploymentModel)
	serviceTemp, err := addServices(deploymentModel)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	cmd.QuestionAndAnswer("what kind of upgrage strategy do you want to use (rolling/recreate) :", func(answer string) {
		if "rolling" == strings.ToLower(answer) {
			deploymentModel.Spec.Strategy = model.DeploymentStrategy{
				Type: "Rolling",
				RollingParams: &model.RollingDeploymentStrategyParams{ //prob need to prompt for these
					UpdatePeriodSeconds: &[]int64{1}[0], //bit crap but it want a pointer rather than value.
					IntervalSeconds:     &[]int64{1}[0],
					TimeoutSeconds:      &[]int64{300}[0],
				},
			}
		}
	})

	if err := templateServ.SaveDeployment(temp, name, deploymentModel); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	if err := templateServ.SaveService(temp, name, serviceTemp); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}
