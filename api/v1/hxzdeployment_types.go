/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HxzDeploymentSpec defines the desired state of HxzDeployment
type HxzDeploymentSpec struct {
	// Image 存储镜像地址
	Image string `json:"image"`
	// Port 存储服务提供的端口
	Port int32 `json:"port"`
	// Replicas 存储要部署多少个副本
	// +optional
	Replicas int32 `json:"replicas,omitempty"`
	// StartCmd 存储启动命令
	// +optional
	StartCmd string `json:"startCmd,omitempty"`
	// Args 存储启动命令参数
	// +optional
	Args []string `json:"args,omitempty"`
	// Environments 存储环境变量，直接使用pod中的定义方式
	// +optional
	Environments []corev1.EnvVar `json:"environments,omitempty"`
	// Expose service要暴露的端口
	Expose *Expose `json:"expose"`
}

// Expose defines the desired state of Expose
type Expose struct {
	// Mode 模式 nodePort or ingress
	Mode string `json:"mode"`
	// IngressDomain 域名。 在 Mode 为 ingress 的时候， 此项为必填
	// +optional
	IngressDomain string `json:"ingressDomain,omitempty"`
	// NodePort nodePort端口。 在 mode 为 nodePort 的时候， 此项为必填
	// +optional
	NodePort int32 `json:"nodePort,omitempty"`
	// ServicePort service的端口，一般是随机生成，这里我们为了防止冲突，使用和我们提供服务相同的端口。
	// +optional
	ServicePort int32 `json:"servicePort,omitempty"`
}

// HxzDeploymentStatus defines the observed state of HxzDeployment
type HxzDeploymentStatus struct {
	// Phase 处于什么阶段
	// +optional
	Phase string `json:"phase,omitempty"`
	// Message 这个阶段的信息
	// +optional
	Message string `json:"message,omitempty"`
	// Reason 处于这个阶段的原因
	// +optional
	Reason string `json:"reason,omitempty"`
	// Conditions 这个阶段的子资源的状态
	// +optional
	Conditions []Condition `json:"conditions,omitempty"`
}

// Condition defines the observed state of Condition
type Condition struct {
	// Type 子资源类型
	// +optional
	Type string `json:"type,omitempty"`
	// Message 这个这个子资源状态的信息
	// +optional
	Message string `json:"message,omitempty"`
	// Status 这个子资源的状态名称
	// +optional
	Status string `json:"status,omitempty"`
	// Reason 处于这个状态的原因
	// +optional
	Reason string `json:"reason,omitempty"`
	// LastTransitionTime 最后创建/更新的时间
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HxzDeployment is the Schema for the hxzdeployments API
type HxzDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HxzDeploymentSpec   `json:"spec,omitempty"`
	Status HxzDeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HxzDeploymentList contains a list of HxzDeployment
type HxzDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HxzDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HxzDeployment{}, &HxzDeploymentList{})
}
