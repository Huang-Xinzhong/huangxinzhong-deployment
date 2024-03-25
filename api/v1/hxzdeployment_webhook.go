/*
Copyright 2022.

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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"strings"
)

// log is for logging in this package.
var hxzdeploymentlog = logf.Log.WithName("hxzdeployment-resource")

func (r *HxzDeployment) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-apps-huangxinzhong-com-v1-hxzdeployment,mutating=true,failurePolicy=fail,sideEffects=None,groups=apps.huangxinzhong.com,resources=hxzdeployments,verbs=create;update,versions=v1,name=mhxzdeployment.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &HxzDeployment{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *HxzDeployment) Default() {
	hxzdeploymentlog.Info("default", "name", r.Name)

	if r.Spec.Replicas == 0 {
		// 因为我们不能确定用户给定的服务是否是一个无状态的应用，如果是有状态的，
		// 多个副本会造成数据错乱。所以我们保守的，只给一个副本
		r.Spec.Replicas = 1
	}

	if r.Spec.Expose.ServicePort == 0 {
		// 允许用户自己指定service的port值，如果不指定，则使用服务的port值来代替
		r.Spec.Expose.ServicePort = r.Spec.Port
	}

	// 增加每个字符串字段的空格处理
}

//+kubebuilder:webhook:path=/validate-apps-huangxinzhong-com-v1-hxzdeployment,mutating=false,failurePolicy=fail,sideEffects=None,groups=apps.huangxinzhong.com,resources=hxzdeployments,verbs=create;update,versions=v1,name=vhxzdeployment.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &HxzDeployment{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *HxzDeployment) ValidateCreate() error {
	hxzdeploymentlog.Info("validate create", "name", r.Name)

	return r.validateCreateAndUpdate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *HxzDeployment) ValidateUpdate(_ runtime.Object) error {
	hxzdeploymentlog.Info("validate update", "name", r.Name)

	return r.validateCreateAndUpdate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *HxzDeployment) ValidateDelete() error {
	hxzdeploymentlog.Info("validate delete", "name", r.Name)
	return nil
}

func (r *HxzDeployment) validateCreateAndUpdate() error {
	// 定义错误切片，在后续出现错误的时候，不断的向其中追加，最后合并返回
	errs := field.ErrorList{}
	exposePath := field.NewPath("spec", "expose")

	// 1. 传入的 spec.expose.mode 值是否为 ingress 或 nodeport
	if strings.ToLower(r.Spec.Expose.Mode) != ModeIngress &&
		strings.ToLower(r.Spec.Expose.Mode) != ModeNodePort {
		errs = append(errs,
			field.NotSupported(
				exposePath,
				r.Spec.Expose.Mode,
				[]string{ModeIngress, ModeNodePort}))
	}

	// 2. 如果传入的 spec.expose.mode 是 ingress，那么，spec.expose.ingressDomain 不能为空
	if strings.ToLower(r.Spec.Expose.Mode) == ModeIngress &&
		r.Spec.Expose.IngressDomain == "" {
		errs = append(errs,
			field.Invalid(
				exposePath,
				r.Spec.Expose.Mode,
				"如果`spec.expose.mode` 是 `ingress`，那么，`spec.expose.ingressDomain` 不能为空"))
	}
	// 3. 如果传入的 spec.expose.mode 是 nodeport，那么，spec.expose.nodeport 取值范围是否是 30000-32767
	if strings.ToLower(r.Spec.Expose.Mode) == ModeNodePort &&
		(r.Spec.Expose.NodePort == 0 ||
			r.Spec.Expose.NodePort < 30000 ||
			r.Spec.Expose.NodePort > 32767) {
		errs = append(errs,
			field.Invalid(
				exposePath,
				r.Spec.Expose.Mode,
				"如果 `spec.expose.mode` 是 `nodeport`，那么，`spec.expose.nodeport` 取值范围是否是 30000-32767"))
	}

	return errs.ToAggregate()
}
