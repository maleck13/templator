package model

import (
	"k8s.io/kubernetes/pkg/api/unversioned"
	k8 "k8s.io/kubernetes/pkg/api/v1"
)

//customised types for ost tool

const (
	ReplicationStrategy_EqualToNodes = "#EqualToNodes" //if it is equal to nodes then the replicas are dynamically set == to the number of nodes
	ReplicationStrategy_Single       = "#Single"       //if it is single then it will always set replicas to 1
	DeploymentStrategy_SingleConfig  = "#SingleConfig"
	DeploymentStrategy_PerNodeConfig = "#PerNodeConfig" //dynamically generate a deployment config per node
)

// DeploymentConfig represents a configuration for a single deployment (represented as a
// ReplicationController). It also contains details about changes which resulted in the current
// state of the DeploymentConfig. Each change to the DeploymentConfig which should result in
// a new deployment results in an increment of LatestVersion.
type OSTDeploymentConfig struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	k8.ObjectMeta `json:"metadata,omitempty"`

	// Spec represents a desired deployment state and how to deploy to it.
	Spec OSTDeploymentConfigSpec `json:"spec"`

	// Status represents the current deployment state.
	Status DeploymentConfigStatus `json:"status"`
}

func NewOstDeploymentConfig(name string) *OSTDeploymentConfig {
	deploymentModel := &OSTDeploymentConfig{}
	deploymentModel.ObjectMeta.Name = name
	deploymentModel.Spec = OSTDeploymentConfigSpec{}

	deploymentModel.Kind = "DeploymentConfig"
	deploymentModel.APIVersion = "v1"
	deploymentModel.Spec.Replicas = 1
	deploymentModel.Spec.Triggers = append(deploymentModel.Spec.Triggers, DeploymentTriggerPolicy{
		Type: DeploymentTriggerType("ConfigChange"),
	})

	deploymentModel.Spec.Template = &k8.PodTemplateSpec{}
	deploymentModel.Spec.Template.ObjectMeta.Name = name
	deploymentModel.Spec.Template.Spec.RestartPolicy = "Always"
	deploymentModel.Spec.Template.Spec.DNSPolicy = "ClusterFirst"

	labels := make(map[string]string)
	labels["name"] = name
	deploymentModel.Spec.Template.Labels = labels
	deploymentModel.ObjectMeta.Labels = labels
	selector := make(map[string]string)
	selector["name"] = name
	deploymentModel.Spec.DeploymentConfigSpec.Selector = selector
	return deploymentModel
}

func (osd *OSTDeploymentConfig) GetObjectKind() unversioned.ObjectKind {
	return &osd.TypeMeta
}

type OSTDeploymentConfigSpec struct {
	DeploymentConfigSpec
	// used to indicate how to dynamically set the number of replicas based on the number of nodes
	ReplicaStrategy string `json:"-"`
	// used to indicate how to dynamically build the number of DeploymentConfigs required based on the number of nodes
	DeploymentStrategy string `json:"-"`
}
