package model

import (
	"k8s.io/kubernetes/pkg/api/unversioned"
	k8 "k8s.io/kubernetes/pkg/api/v1"
)

func NewApplicationTemplate(name string) *ApplicationTemplate {
	at := &ApplicationTemplate{}
	at.ObjectMeta = k8.ObjectMeta{}
	at.ObjectMeta.Name = name
	at.ObjectMeta.Annotations = make(map[string]string)
	at.ObjectMeta.Annotations["description"] = "a generated template for " + name
	at.DeploymentConfigs = make(map[string]*OSTDeploymentConfig)
	at.Parameters = make([]*Parameter, 0)
	at.PersistentVolumes = make(map[string]*k8.PersistentVolumeClaim)
	at.Services = make(map[string]*k8.Service)
	at.Pods = make(map[string]*k8.Pod)
	at.Routes = make(map[string]*Route)
	at.ObjectMeta.Name = name
	at.APIVersion = "v1" //todo not hard coded
	at.Kind = "Template"
	return at
}

type ApplicationTemplate struct {
	unversioned.TypeMeta `json:",inline"`
	k8.ObjectMeta        `json:"metadata,omitempty"`
	Services             map[string]*k8.Service               `json:"services"`
	DeploymentConfigs    map[string]*OSTDeploymentConfig      `json:"deploymentConfigs"`
	PersistentVolumes    map[string]*k8.PersistentVolumeClaim `json:"persistentVolumes"`
	Pods                 map[string]*k8.Pod                   `json:"pods"`
	Routes               map[string]*Route                    `json:"routes"`
	Parameters           []*Parameter                         `json:"parameters"`
}
