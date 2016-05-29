package model

import (
	"k8s.io/kubernetes/pkg/api/unversioned"
	k8 "k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/intstr"
)

//******************* copied from openshift origin to save dragging in all of origin as a dep ********************//

// DeploymentPhase describes the possible states a deployment can be in.
type DeploymentPhase string

// DeploymentStrategy describes how to perform a deployment.
type DeploymentStrategy struct {
	// Type is the name of a deployment strategy.
	Type DeploymentStrategyType `json:"type,omitempty"`

	// CustomParams are the input to the Custom deployment strategy.
	CustomParams *CustomDeploymentStrategyParams `json:"customParams,omitempty"`
	// RecreateParams are the input to the Recreate deployment strategy.
	RecreateParams *RecreateDeploymentStrategyParams `json:"recreateParams,omitempty"`
	// RollingParams are the input to the Rolling deployment strategy.
	RollingParams *RollingDeploymentStrategyParams `json:"rollingParams,omitempty"`

	// Resources contains resource requirements to execute the deployment and any hooks
	Resources k8.ResourceRequirements `json:"resources,omitempty"`
	// Labels is a set of key, value pairs added to custom deployer and lifecycle pre/post hook pods.
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is a set of key, value pairs added to custom deployer and lifecycle pre/post hook pods.
	Annotations map[string]string `json:"annotations,omitempty"`
}

// DeploymentStrategyType refers to a specific DeploymentStrategy implementation.
type DeploymentStrategyType string

// CustomDeploymentStrategyParams are the input to the Custom deployment strategy.
type CustomDeploymentStrategyParams struct {
	// Image specifies a Docker image which can carry out a deployment.
	Image string `json:"image,omitempty"`
	// Environment holds the environment which will be given to the container for Image.
	Environment []k8.EnvVar `json:"environment,omitempty"`
	// Command is optional and overrides CMD in the container Image.
	Command []string `json:"command,omitempty"`
}

// RecreateDeploymentStrategyParams are the input to the Recreate deployment
// strategy.
type RecreateDeploymentStrategyParams struct {
	// TimeoutSeconds is the time to wait for updates before giving up. If the
	// value is nil, a default will be used.
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`
	// Pre is a lifecycle hook which is executed before the strategy manipulates
	// the deployment. All LifecycleHookFailurePolicy values are supported.
	Pre *LifecycleHook `json:"pre,omitempty"`
	// Mid is a lifecycle hook which is executed while the deployment is scaled down to zero before the first new
	// pod is created. All LifecycleHookFailurePolicy values are supported.
	Mid *LifecycleHook `json:"mid,omitempty"`
	// Post is a lifecycle hook which is executed after the strategy has
	// finished all deployment logic. All LifecycleHookFailurePolicy values are supported.
	Post *LifecycleHook `json:"post,omitempty"`
}

// LifecycleHook defines a specific deployment lifecycle action. Only one type of action may be specified at any time.
type LifecycleHook struct {
	// FailurePolicy specifies what action to take if the hook fails.
	FailurePolicy LifecycleHookFailurePolicy `json:"failurePolicy"`

	// ExecNewPod specifies the options for a lifecycle hook backed by a pod.
	ExecNewPod *ExecNewPodHook `json:"execNewPod,omitempty"`

	// TagImages instructs the deployer to tag the current image referenced under a container onto an image stream tag.
	TagImages []TagImageHook `json:"tagImages,omitempty"`
}

// LifecycleHookFailurePolicy describes possibles actions to take if a hook fails.
type LifecycleHookFailurePolicy string

// ExecNewPodHook is a hook implementation which runs a command in a new pod
// based on the specified container which is assumed to be part of the
// deployment template.
type ExecNewPodHook struct {
	// Command is the action command and its arguments.
	Command []string `json:"command"`
	// Env is a set of environment variables to supply to the hook pod's container.
	Env []k8.EnvVar `json:"env,omitempty"`
	// ContainerName is the name of a container in the deployment pod template
	// whose Docker image will be used for the hook pod's container.
	ContainerName string `json:"containerName"`
	// Volumes is a list of named volumes from the pod template which should be
	// copied to the hook pod. Volumes names not found in pod spec are ignored.
	// An empty list means no volumes will be copied.
	Volumes []string `json:"volumes,omitempty"`
}

// TagImageHook is a request to tag the image in a particular container onto an ImageStreamTag.
type TagImageHook struct {
	// ContainerName is the name of a container in the deployment config whose image value will be used as the source of the tag. If there is only a single
	// container this value will be defaulted to the name of that container.
	ContainerName string `json:"containerName"`
	// To is the target ImageStreamTag to set the container's image onto.
	To k8.ObjectReference `json:"to"`
}

// RollingDeploymentStrategyParams are the input to the Rolling deployment
// strategy.
type RollingDeploymentStrategyParams struct {
	// UpdatePeriodSeconds is the time to wait between individual pod updates.
	// If the value is nil, a default will be used.
	UpdatePeriodSeconds *int64 `json:"updatePeriodSeconds,omitempty"`
	// IntervalSeconds is the time to wait between polling deployment status
	// after update. If the value is nil, a default will be used.
	IntervalSeconds *int64 `json:"intervalSeconds,omitempty"`
	// TimeoutSeconds is the time to wait for updates before giving up. If the
	// value is nil, a default will be used.
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`
	// MaxUnavailable is the maximum number of pods that can be unavailable
	// during the update. Value can be an absolute number (ex: 5) or a
	// percentage of total pods at the start of update (ex: 10%). Absolute
	// number is calculated from percentage by rounding up.
	//
	// This cannot be 0 if MaxSurge is 0. By default, 25% is used.
	//
	// Example: when this is set to 30%, the old RC can be scaled down by 30%
	// immediately when the rolling update starts. Once new pods are ready, old
	// RC can be scaled down further, followed by scaling up the new RC,
	// ensuring that at least 70% of original number of pods are available at
	// all times during the update.
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`
	// MaxSurge is the maximum number of pods that can be scheduled above the
	// original number of pods. Value can be an absolute number (ex: 5) or a
	// percentage of total pods at the start of the update (ex: 10%). Absolute
	// number is calculated from percentage by rounding up.
	//
	// This cannot be 0 if MaxUnavailable is 0. By default, 25% is used.
	//
	// Example: when this is set to 30%, the new RC can be scaled up by 30%
	// immediately when the rolling update starts. Once old pods have been
	// killed, new RC can be scaled up further, ensuring that total number of
	// pods running at any time during the update is atmost 130% of original
	// pods.
	MaxSurge *intstr.IntOrString `json:"maxSurge,omitempty"`
	// UpdatePercent is the percentage of replicas to scale up or down each
	// interval. If nil, one replica will be scaled up and down each interval.
	// If negative, the scale order will be down/up instead of up/down.
	// DEPRECATED: Use MaxUnavailable/MaxSurge instead.
	UpdatePercent *int `json:"updatePercent,omitempty"`
	// Pre is a lifecycle hook which is executed before the deployment process
	// begins. All LifecycleHookFailurePolicy values are supported.
	Pre *LifecycleHook `json:"pre,omitempty"`
	// Post is a lifecycle hook which is executed after the strategy has
	// finished all deployment logic. The LifecycleHookFailurePolicyAbort policy
	// is NOT supported.
	Post *LifecycleHook `json:"post,omitempty"`
}

// DeploymentConfig represents a configuration for a single deployment (represented as a
// ReplicationController). It also contains details about changes which resulted in the current
// state of the DeploymentConfig. Each change to the DeploymentConfig which should result in
// a new deployment results in an increment of LatestVersion.
type DeploymentConfig struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	k8.ObjectMeta `json:"metadata,omitempty"`

	// Spec represents a desired deployment state and how to deploy to it.
	Spec DeploymentConfigSpec `json:"spec"`

	// Status represents the current deployment state.
	Status DeploymentConfigStatus `json:"status"`
}

// DeploymentConfigSpec represents the desired state of the deployment.
type DeploymentConfigSpec struct {
	// Strategy describes how a deployment is executed.
	Strategy DeploymentStrategy `json:"strategy"`

	// Triggers determine how updates to a DeploymentConfig result in new deployments. If no triggers
	// are defined, a new deployment can only occur as a result of an explicit client update to the
	// DeploymentConfig with a new LatestVersion.
	Triggers []DeploymentTriggerPolicy `json:"triggers"`

	// Replicas is the number of desired replicas.
	Replicas int `json:"replicas"`

	// Test ensures that this deployment config will have zero replicas except while a deployment is running. This allows the
	// deployment config to be used as a continuous deployment test - triggering on images, running the deployment, and then succeeding
	// or failing. Post strategy hooks and After actions can be used to integrate successful deployment with an action.
	Test bool `json:"test"`

	// Selector is a label query over pods that should match the Replicas count.
	Selector map[string]string `json:"selector,omitempty"`

	// Template is the object that describes the pod that will be created if
	// insufficient replicas are detected.
	Template *k8.PodTemplateSpec `json:"template,omitempty"`
}

// DeploymentConfigStatus represents the current deployment state.
type DeploymentConfigStatus struct {
	// LatestVersion is used to determine whether the current deployment associated with a DeploymentConfig
	// is out of sync.
	LatestVersion int `json:"latestVersion,omitempty"`
	// Details are the reasons for the update to this deployment config.
	// This could be based on a change made by the user or caused by an automatic trigger
	Details *DeploymentDetails `json:"details,omitempty"`
	// ObservedGeneration is the most recent generation observed by the controller.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// DeploymentTriggerPolicy describes a policy for a single trigger that results in a new deployment.
type DeploymentTriggerPolicy struct {
	// Type of the trigger
	Type DeploymentTriggerType `json:"type,omitempty"`
	// ImageChangeParams represents the parameters for the ImageChange trigger.
	ImageChangeParams *DeploymentTriggerImageChangeParams `json:"imageChangeParams,omitempty"`
}

// DeploymentTriggerType refers to a specific DeploymentTriggerPolicy implementation.
type DeploymentTriggerType string

// DeploymentTriggerImageChangeParams represents the parameters to the ImageChange trigger.
type DeploymentTriggerImageChangeParams struct {
	// Automatic means that the detection of a new tag value should result in an image update
	// inside the pod template. Deployment configs that haven't been deployed yet will always
	// have their images updated. Deployment configs that have been deployed at least once, will
	// have their images updated only if this is set to true.
	Automatic bool `json:"automatic,omitempty"`
	// ContainerNames is used to restrict tag updates to the specified set of container names in a pod.
	ContainerNames []string `json:"containerNames,omitempty"`
	// From is a reference to an image stream tag to watch for changes. From.Name is the only
	// required subfield - if From.Namespace is blank, the namespace of the current deployment
	// trigger will be used.
	From k8.ObjectReference `json:"from"`
	// LastTriggeredImage is the last image to be triggered.
	LastTriggeredImage string `json:"lastTriggeredImage,omitempty"`
}

// DeploymentDetails captures information about the causes of a deployment.
type DeploymentDetails struct {
	// Message is the user specified change message, if this deployment was triggered manually by the user
	Message string `json:"message,omitempty"`
	// Causes are extended data associated with all the causes for creating a new deployment
	Causes []DeploymentCause `json:"causes"`
}

// DeploymentCause captures information about a particular cause of a deployment.
type DeploymentCause struct {
	// Type of the trigger that resulted in the creation of a new deployment
	Type DeploymentTriggerType `json:"type"`
	// ImageTrigger contains the image trigger details, if this trigger was fired based on an image change
	ImageTrigger *DeploymentCauseImageTrigger `json:"imageTrigger,omitempty"`
}

// DeploymentCauseImageTrigger represents details about the cause of a deployment originating
// from an image change trigger
type DeploymentCauseImageTrigger struct {
	// From is a reference to the changed object which triggered a deployment. The field may have
	// the kinds DockerImage, ImageStreamTag, or ImageStreamImage.
	From k8.ObjectReference `json:"from"`
}

// DeploymentConfigList is a collection of deployment configs.
type DeploymentConfigList struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	unversioned.ListMeta `json:"metadata,omitempty"`

	// Items is a list of deployment configs
	Items []DeploymentConfig `json:"items"`
}

// DeploymentConfigRollback provides the input to rollback generation.
type DeploymentConfigRollback struct {
	unversioned.TypeMeta `json:",inline"`
	// Spec defines the options to rollback generation.
	Spec DeploymentConfigRollbackSpec `json:"spec"`
}

// DeploymentConfigRollbackSpec represents the options for rollback generation.
type DeploymentConfigRollbackSpec struct {
	// From points to a ReplicationController which is a deployment.
	From k8.ObjectReference `json:"from"`
	// IncludeTriggers specifies whether to include config Triggers.
	IncludeTriggers bool `json:"includeTriggers"`
	// IncludeTemplate specifies whether to include the PodTemplateSpec.
	IncludeTemplate bool `json:"includeTemplate"`
	// IncludeReplicationMeta specifies whether to include the replica count and selector.
	IncludeReplicationMeta bool `json:"includeReplicationMeta"`
	// IncludeStrategy specifies whether to include the deployment Strategy.
	IncludeStrategy bool `json:"includeStrategy"`
}

// DeploymentLog represents the logs for a deployment
type DeploymentLog struct {
	unversioned.TypeMeta `json:",inline"`
}

// DeploymentLogOptions is the REST options for a deployment log
type DeploymentLogOptions struct {
	unversioned.TypeMeta `json:",inline"`

	// The container for which to stream logs. Defaults to only container if there is one container in the pod.
	Container string `json:"container,omitempty"`
	// Follow if true indicates that the build log should be streamed until
	// the build terminates.
	Follow bool `json:"follow,omitempty"`
	// Return previous deployment logs. Defaults to false.
	Previous bool `json:"previous,omitempty"`
	// A relative time in seconds before the current time from which to show logs. If this value
	// precedes the time a pod was started, only logs since the pod start will be returned.
	// If this value is in the future, no logs will be returned.
	// Only one of sinceSeconds or sinceTime may be specified.
	SinceSeconds *int64 `json:"sinceSeconds,omitempty"`
	// An RFC3339 timestamp from which to show logs. If this value
	// preceeds the time a pod was started, only logs since the pod start will be returned.
	// If this value is in the future, no logs will be returned.
	// Only one of sinceSeconds or sinceTime may be specified.
	SinceTime *unversioned.Time `json:"sinceTime,omitempty"`
	// If true, add an RFC3339 or RFC3339Nano timestamp at the beginning of every line
	// of log output. Defaults to false.
	Timestamps bool `json:"timestamps,omitempty"`
	// If set, the number of lines from the end of the logs to show. If not specified,
	// logs are shown from the creation of the container or sinceSeconds or sinceTime
	TailLines *int64 `json:"tailLines,omitempty"`
	// If set, the number of bytes to read from the server before terminating the
	// log output. This may not display a complete final line of logging, and may return
	// slightly more or slightly less than the specified limit.
	LimitBytes *int64 `json:"limitBytes,omitempty"`

	// NoWait if true causes the call to return immediately even if the deployment
	// is not available yet. Otherwise the server will wait until the deployment has started.
	// TODO: Fix the tag to 'noWait' in v2
	NoWait bool `json:"nowait,omitempty"`

	// Version of the deployment for which to view logs.
	Version *int64 `json:"version,omitempty"`
}

// Parameter defines a name/value variable that is to be processed during
// the Template to Config transformation.
type Parameter struct {
	// Required: Parameter name must be set and it can be referenced in Template
	// Items using ${PARAMETER_NAME}
	Name string

	// Optional: The name that will show in UI instead of parameter 'Name'
	DisplayName string

	// Optional: Parameter can have description
	Description string

	// Optional: Value holds the Parameter data. If specified, the generator
	// will be ignored. The value replaces all occurrences of the Parameter
	// ${Name} expression during the Template to Config transformation.
	Value string

	// Optional: Generate specifies the generator to be used to generate
	// random string from an input value specified by From field. The result
	// string is stored into Value field. If empty, no generator is being
	// used, leaving the result Value untouched.
	Generate string

	// Optional: From is an input value for the generator.
	From string

	// Optional: Indicates the parameter must have a value.  Defaults to false.
	Required bool
}

// Route encapsulates the inputs needed to connect an alias to endpoints.
type Route struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	k8.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the route
	Spec RouteSpec `json:"spec"`
	// Status is the current state of the route
	Status RouteStatus `json:"status"`
}

func (r *Route) GetObjectKind() unversioned.ObjectKind {
	return &r.TypeMeta
}

// RouteList is a collection of Routes.
type RouteList struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	unversioned.ListMeta `json:"metadata,omitempty"`

	// Items is a list of routes
	Items []Route `json:"items"`
}

// RouteSpec describes the route the user wishes to exist.
type RouteSpec struct {
	// Ports are the ports that the user wishes to expose.
	//Ports []RoutePort `json:"ports,omitempty"`

	// Host is an alias/DNS that points to the service. Optional
	// Must follow DNS952 subdomain conventions.
	Host string `json:"host"`
	// Path that the router watches for, to route traffic for to the service. Optional
	Path string `json:"path,omitempty"`

	// To is an object the route points to. Only the Service kind is allowed, and it will
	// be defaulted to Service.
	To k8.ObjectReference `json:"to"`

	// If specified, the port to be used by the router. Most routers will use all
	// endpoints exposed by the service by default - set this value to instruct routers
	// which port to use.
	Port *RoutePort `json:"port,omitempty"`

	// TLS provides the ability to configure certificates and termination for the route
	TLS *TLSConfig `json:"tls,omitempty"`
}

// RoutePort defines a port mapping from a router to an endpoint in the service endpoints.
type RoutePort struct {
	// The target port on pods selected by the service this route points to.
	// If this is a string, it will be looked up as a named port in the target
	// endpoints port list. Required
	TargetPort intstr.IntOrString `json:"targetPort"`
}

// RouteStatus provides relevant info about the status of a route, including which routers
// acknowledge it.
type RouteStatus struct {
	// Ingress describes the places where the route may be exposed. The list of
	// ingress points may contain duplicate Host or RouterName values. Routes
	// are considered live once they are `Ready`
	Ingress []RouteIngress `json:"ingress"`
}

// RouteIngress holds information about the places where a route is exposed
type RouteIngress struct {
	// Host is the host string under which the route is exposed; this value is required
	Host string `json:"host,omitempty"`
	// Name is a name chosen by the router to identify itself; this value is required
	RouterName string `json:"routerName,omitempty"`
	// Conditions is the state of the route, may be empty.
	Conditions []RouteIngressCondition `json:"conditions,omitempty"`
}

// RouteIngressCondition contains details for the current condition of this pod.
// TODO: add LastTransitionTime, Reason, Message to match NodeCondition api.
type RouteIngressCondition struct {
	// Type is the type of the condition.
	// Currently only Ready.
	Type RouteIngressConditionType `json:"type"`
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	Status k8.ConditionStatus `json:"status"`
	// (brief) reason for the condition's last transition, and is usually a machine and human
	// readable constant
	Reason string `json:"reason,omitempty"`
	// Human readable message indicating details about last transition.
	Message string `json:"message,omitempty"`
	// RFC 3339 date and time when this condition last transitioned
	LastTransitionTime *unversioned.Time `json:"lastTransitionTime,omitempty"`
}

// RouteIngressConditionType is a valid value for RouteCondition
type RouteIngressConditionType string

// RouterShard has information of a routing shard and is used to
// generate host names and routing table entries when a routing shard is
// allocated for a specific route.
// Caveat: This is WIP and will likely undergo modifications when sharding
//         support is added.
type RouterShard struct {
	// ShardName uniquely identifies a router shard in the "set" of
	// routers used for routing traffic to the services.
	ShardName string `json:"shardName"`

	// DNSSuffix for the shard ala: shard-1.v3.openshift.com
	DNSSuffix string `json:"dnsSuffix"`
}

// TLSConfig defines config used to secure a route and provide termination
type TLSConfig struct {
	// Termination indicates termination type.
	Termination TLSTerminationType `json:"termination"`

	// Certificate provides certificate contents
	Certificate string `json:"certificate,omitempty"`

	// Key provides key file contents
	Key string `json:"key,omitempty"`

	// CACertificate provides the cert authority certificate contents
	CACertificate string `json:"caCertificate,omitempty"`

	// DestinationCACertificate provides the contents of the ca certificate of the final destination.  When using reencrypt
	// termination this file should be provided in order to have routers use it for health checks on the secure connection
	DestinationCACertificate string `json:"destinationCACertificate,omitempty"`

	// InsecureEdgeTerminationPolicy indicates the desired behavior for
	// insecure connections to an edge-terminated route:
	//   disable, allow or redirect
	InsecureEdgeTerminationPolicy InsecureEdgeTerminationPolicyType `json:"insecureEdgeTerminationPolicy,omitempty"`
}

// TLSTerminationType dictates where the secure communication will stop
// TODO: Reconsider this type in v2
type TLSTerminationType string

// InsecureEdgeTerminationPolicyType dictates the behavior of insecure
// connections to an edge-terminated route.
type InsecureEdgeTerminationPolicyType string

// Template contains the inputs needed to produce a Config.
type Template struct {
	unversioned.TypeMeta
	k8.ObjectMeta `json:"metadata"`

	// Optional: Parameters is an array of Parameters used during the
	// Template to Config transformation.
	Parameters []Parameter `json:"parameters"`

	// Required: A list of resources to create
	Objects []runtime.Object `json:"objects"`

	// Optional: ObjectLabels is a set of labels that are applied to every
	// object during the Template to Config transformation
	ObjectLabels map[string]string `json:"objectLabels"`
}
