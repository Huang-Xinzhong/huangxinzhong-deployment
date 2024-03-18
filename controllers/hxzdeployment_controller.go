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

package controllers

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
	"time"

	myAppsv1 "github.com/huangxinzhong/huangxinzhong-deployment/api/v1"
)

var WaitRqeueue = 10 * time.Second

// HxzDeploymentReconciler reconciles a HxzDeployment object
type HxzDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.huangxinzhong.com,resources=hxzdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.huangxinzhong.com,resources=hxzdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.huangxinzhong.com,resources=hxzdeployments/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HxzDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *HxzDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// 状态更新策略
	// 创建的时候
	//     更新为创建
	// 更新的时候
	// 	   根据获取的状态来判断时候更新 status
	// 删除的时候
	// 	  只有在操作 ingress 的时候，并且 mode 为 nodeport 的时候

	logger := log.FromContext(ctx, "HxzDeployment", req.NamespacedName)

	logger.Info("Reconcile is started.")
	// 1. 获取资源对象
	md := new(myAppsv1.HxzDeployment)
	if err := r.Client.Get(ctx, req.NamespacedName, md); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 防止污染缓存
	mdCopy := md.DeepCopy()

	// ====== 处理 deployment=======
	// 2. 获取 deployment 资源对象
	deploy := new(appsv1.Deployment)
	if err := r.Client.Get(ctx, req.NamespacedName, deploy); err != nil {
		if errors.IsNotFound(err) {
			// 2.1 不存在对象
			// 2.1.1 创建 deployment
			if errCreate := r.createDeployment(ctx, mdCopy); errCreate != nil {
				return ctrl.Result{}, errCreate
			}

			if _, errStatus := r.updateStatus(
				ctx,
				mdCopy,
				myAppsv1.ConditionTypeDeployment,
				fmt.Sprintf(myAppsv1.ConditionMessageDeploymentNotFmt, req.Name),
				myAppsv1.ConditionStatusFalse,
				myAppsv1.ConditionReasonDeploymentNotReady,
			); errStatus != nil {
				return ctrl.Result{}, errStatus
			}

		} else {
			if _, errStatus := r.updateStatus(
				ctx,
				mdCopy,
				myAppsv1.ConditionTypeDeployment,
				fmt.Sprintf("Deployment update failed, err: %s", err.Error()),
				myAppsv1.ConditionStatusFalse,
				myAppsv1.ConditionReasonDeploymentNotReady,
			); errStatus != nil {
				return ctrl.Result{}, errStatus
			}

			return ctrl.Result{}, err
		}
	} else {
		// 2.2 存在对象
		// 2.2.1 更新 deployment
		if err := r.updateDeployment(ctx, mdCopy, deploy); err != nil {
			return ctrl.Result{}, err
		}

		if deploy.Status.AvailableReplicas == mdCopy.Spec.Replicas {
			if _, errStatus := r.updateStatus(
				ctx,
				mdCopy,
				myAppsv1.ConditionTypeDeployment,
				fmt.Sprintf(myAppsv1.ConditionMessageDeploymentOkFmt, req.Name),
				myAppsv1.ConditionStatusTrue,
				myAppsv1.ConditionReasonDeploymentReady,
			); errStatus != nil {
				return ctrl.Result{}, errStatus
			}
		} else {
			if _, errStatus := r.updateStatus(
				ctx,
				mdCopy,
				myAppsv1.ConditionTypeDeployment,
				fmt.Sprintf(myAppsv1.ConditionMessageDeploymentNotFmt, req.Name),
				myAppsv1.ConditionStatusFalse,
				myAppsv1.ConditionReasonDeploymentNotReady,
			); errStatus != nil {
				return ctrl.Result{}, errStatus
			}
		}

	}

	// ====== 处理 Service =========
	// 3. 获取 service 资源对象
	svc := new(corev1.Service)
	if err := r.Client.Get(ctx, req.NamespacedName, svc); err != nil {
		if errors.IsNotFound(err) {
			// 3.1 不存在
			// 3.1.1 mode 为 ingress
			mode := strings.ToLower(mdCopy.Spec.Expose.Mode)
			if mode == myAppsv1.ModeIngress {
				// 3.1.1.1 创建普通 service
				if err := r.createService(ctx, mdCopy); err != nil {
					return ctrl.Result{}, err
				}
			} else if mode == myAppsv1.ModeNodePort {
				// 3.1.2 mode 为 nodeport
				// 3.1.2.1 创建 nodeport 模式的 service
				if err := r.createNodePortService(ctx, mdCopy); err != nil {
					return ctrl.Result{}, err
				}
			} else {
				return ctrl.Result{}, myAppsv1.ErrorNotSupportMode
			}

			if _, errStatus := r.updateStatus(
				ctx,
				mdCopy,
				myAppsv1.ConditionTypeService,
				fmt.Sprintf(myAppsv1.ConditionMessageServiceNotFmt, req.Name),
				myAppsv1.ConditionStatusFalse,
				myAppsv1.ConditionReasonServiceNotReady,
			); errStatus != nil {
				return ctrl.Result{}, errStatus
			}
		} else {
			if _, errStatus := r.updateStatus(
				ctx,
				mdCopy,
				myAppsv1.ConditionTypeService,
				fmt.Sprintf("Service update failed, err: %s", err.Error()),
				myAppsv1.ConditionStatusFalse,
				myAppsv1.ConditionReasonServiceNotReady,
			); errStatus != nil {
				return ctrl.Result{}, errStatus
			}
			return ctrl.Result{}, err
		}
	} else {
		// 3.2 存在
		mode := strings.ToLower(mdCopy.Spec.Expose.Mode)
		if mode == myAppsv1.ModeIngress {
			// 3.2.1.1 更新普通 service
			if err := r.updateService(ctx, mdCopy, svc); err != nil {
				return ctrl.Result{}, err
			}
		} else if mode == myAppsv1.ModeNodePort {
			// 3.2.2 mode 为 nodeport
			// 3.2.2.1 更新 nodeport 模式的 service
			if err := r.updateNodePortService(ctx, mdCopy, svc); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			return ctrl.Result{}, myAppsv1.ErrorNotSupportMode
		}

		if _, errStatus := r.updateStatus(
			ctx,
			mdCopy,
			myAppsv1.ConditionTypeService,
			fmt.Sprintf(myAppsv1.ConditionMessageServiceOkFmt, req.Name),
			myAppsv1.ConditionStatusTrue,
			myAppsv1.ConditionReasonServiceReady,
		); errStatus != nil {
			return ctrl.Result{}, errStatus
		}
	}

	// ====== 处理 ingress =========
	// 4 获取 ingress 资源对象
	ig := new(networkv1.Ingress)
	if err := r.Client.Get(ctx, req.NamespacedName, ig); err != nil {
		if errors.IsNotFound(err) {
			// 4.1 不存在

			mode := strings.ToLower(mdCopy.Spec.Expose.Mode)
			if mode == myAppsv1.ModeIngress {
				// 4.1.1 mode 为 ingress
				// 4.1.1.1 创建 ingress
				if err := r.createIngress(ctx, mdCopy); err != nil {
					return ctrl.Result{}, err
				}

				if _, errStatus := r.updateStatus(
					ctx,
					mdCopy,
					myAppsv1.ConditionTypeIngress,
					fmt.Sprintf(myAppsv1.ConditionMessageIngressNotFmt, req.Name),
					myAppsv1.ConditionStatusFalse,
					myAppsv1.ConditionReasonIngressNotReady,
				); errStatus != nil {
					return ctrl.Result{}, errStatus
				}

			} else if mode == myAppsv1.ModeNodePort {
				// 4.1.2 mode 为 nodeport
				// 4.1.2.1 退出
				return ctrl.Result{}, nil
			}
		} else {
			if _, errStatus := r.updateStatus(
				ctx,
				mdCopy,
				myAppsv1.ConditionTypeIngress,
				fmt.Sprintf("Ingress  update failed, err: %s", err.Error()),
				myAppsv1.ConditionStatusFalse,
				myAppsv1.ConditionReasonIngressNotReady,
			); errStatus != nil {
				return ctrl.Result{}, errStatus
			}
			return ctrl.Result{}, err
		}
	} else {
		// 4.2 存在
		mode := strings.ToLower(mdCopy.Spec.Expose.Mode)
		if mode == myAppsv1.ModeIngress {
			// 4.2.1 mode 为 ingress
			// 4.2.1.1 更新 ingress
			if err := r.updateIngress(ctx, mdCopy, ig); err != nil {
				return ctrl.Result{}, err
			}

			if _, errStatus := r.updateStatus(
				ctx,
				mdCopy,
				myAppsv1.ConditionTypeIngress,
				fmt.Sprintf(myAppsv1.ConditionMessageIngressOkFmt, req.Name),
				myAppsv1.ConditionStatusTrue,
				myAppsv1.ConditionReasonIngressReady,
			); errStatus != nil {
				return ctrl.Result{}, errStatus
			}
		} else if mode == myAppsv1.ModeNodePort {
			// 4.2.2 mode 为 nodeport
			// 4.2.2.1 删除 ingress
			if err := r.deleteIngress(ctx, mdCopy); err != nil {
				return ctrl.Result{}, err
			}
			r.deleteStatus(ctx, mdCopy, myAppsv1.ConditionTypeIngress)
		}

	}

	// 最后检查状态是否最终完成
	if success, errStatus := r.updateStatus(
		ctx,
		mdCopy,
		"",
		"",
		"",
		"",
	); errStatus != nil {
		return ctrl.Result{}, errStatus
	} else if !success {
		logger.Info("Reconcile is ended.")
		return ctrl.Result{RequeueAfter: WaitRqeueue}, nil
	}

	logger.Info("Reconcile is ended.")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HxzDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myAppsv1.HxzDeployment{}).
		Owns(&appsv1.Deployment{}). // 监控 deploy 类型， 变更就触发 reconciler
		Owns(&corev1.Service{}).    // 监控 service 类型， 变更就触发 reconciler
		Owns(&networkv1.Ingress{}). // 监控 ingress 类型， 变更就触发 reconciler
		Complete(r)
}

func (r *HxzDeploymentReconciler) createDeployment(ctx context.Context, md *myAppsv1.HxzDeployment) error {
	deploy, err := NewDeployment(md)
	if err != nil {
		return err
	}

	// 设置 deployment 所属于 md
	if err := controllerutil.SetControllerReference(md, deploy, r.Scheme); err != nil {
		return err
	}

	return r.Client.Create(ctx, deploy)
}

func (r *HxzDeploymentReconciler) updateDeployment(ctx context.Context, md *myAppsv1.HxzDeployment, dp *appsv1.Deployment) error {
	deploy, err := NewDeployment(md)
	if err != nil {
		return err
	}

	// 设置 deployment 所属于 md
	if err := controllerutil.SetControllerReference(md, deploy, r.Scheme); err != nil {
		return err
	}

	// 预更新 deployment， 得到更新后的数据
	if err := r.Update(ctx, deploy, &client.DryRunAll); err != nil {
		return err
	}

	// 和之前的数据进行比较， 如果相同， 说明不需要更新
	if reflect.DeepEqual(dp.Spec, deploy.Spec) {
		return nil
	}

	return r.Client.Update(ctx, deploy)
}

func (r *HxzDeploymentReconciler) createService(ctx context.Context, md *myAppsv1.HxzDeployment) error {
	svc, err := NewService(md)
	if err != nil {
		return err
	}

	// 设置 deployment 所属于 md
	if err := controllerutil.SetControllerReference(md, svc, r.Scheme); err != nil {
		return err
	}

	return r.Client.Create(ctx, svc)
}

func (r *HxzDeploymentReconciler) createNodePortService(ctx context.Context, md *myAppsv1.HxzDeployment) error {
	svc, err := NewServiceNodePort(md)
	if err != nil {
		return err
	}

	// 设置 deployment 所属于 md
	if err := controllerutil.SetControllerReference(md, svc, r.Scheme); err != nil {
		return err
	}

	return r.Client.Create(ctx, svc)
}

func (r *HxzDeploymentReconciler) updateService(ctx context.Context, md *myAppsv1.HxzDeployment, svc *corev1.Service) error {
	service, err := NewService(md)
	if err != nil {
		return err
	}

	// 设置 svc 所属于 md
	if err := controllerutil.SetControllerReference(md, service, r.Scheme); err != nil {
		return err
	}

	// 预更新 svc， 得到更新后的数据
	if err := r.Update(ctx, service, &client.DryRunAll); err != nil {
		return err
	}

	// 和之前的数据进行比较， 如果相同， 说明不需要更新
	if reflect.DeepEqual(svc.Spec, service.Spec) {
		return nil
	}

	return r.Client.Update(ctx, service)
}

func (r *HxzDeploymentReconciler) updateNodePortService(ctx context.Context, md *myAppsv1.HxzDeployment, svc *corev1.Service) error {
	service, err := NewServiceNodePort(md)
	if err != nil {
		return err
	}

	// 设置 deployment 所属于 md
	if err := controllerutil.SetControllerReference(md, service, r.Scheme); err != nil {
		return err
	}

	// 预更新 svc， 得到更新后的数据
	if err := r.Update(ctx, service, &client.DryRunAll); err != nil {
		return err
	}

	// 和之前的数据进行比较， 如果相同， 说明不需要更新
	if reflect.DeepEqual(svc.Spec, service.Spec) {
		return nil
	}
	return r.Client.Update(ctx, service)
}

func (r *HxzDeploymentReconciler) createIngress(ctx context.Context, md *myAppsv1.HxzDeployment) error {
	ingress, err := NewIngress(md)
	if err != nil {
		return err
	}

	// 设置 deployment 所属于 md
	if err := controllerutil.SetControllerReference(md, ingress, r.Scheme); err != nil {
		return err
	}

	return r.Client.Create(ctx, ingress)
}

func (r *HxzDeploymentReconciler) updateIngress(ctx context.Context, md *myAppsv1.HxzDeployment, ig *networkv1.Ingress) error {
	ingress, err := NewIngress(md)
	if err != nil {
		return err
	}

	// 设置 deployment 所属于 md
	if err := controllerutil.SetControllerReference(md, ingress, r.Scheme); err != nil {
		return err
	}

	// 预更新 svc， 得到更新后的数据
	if err := r.Update(ctx, ingress, &client.DryRunAll); err != nil {
		return err
	}

	// 和之前的数据进行比较， 如果相同， 说明不需要更新
	if reflect.DeepEqual(ig.Spec, ingress.Spec) {
		return nil
	}

	return r.Client.Update(ctx, ingress)
}

func (r *HxzDeploymentReconciler) deleteIngress(ctx context.Context, md *myAppsv1.HxzDeployment) error {
	ingress, err := NewIngress(md)
	if err != nil {
		return err
	}

	return r.Client.Delete(ctx, ingress)
}

// 处理 status
// return:
//   - bool: 资源是否完成，是否需要等待，如果是 true， 表示资源已经完成不需要再次 reconcile
//     如果是 false，表示资源还未完成， 需要重新入队
//
// error： 执行 update 的状态
func (r *HxzDeploymentReconciler) updateStatus(ctx context.Context, md *myAppsv1.HxzDeployment, conditionType, message, status, reason string) (bool, error) {
	if conditionType != "" {
		// 1. 获取 status
		// 2. 获取 condition 字段
		// 3. 根据当前的需求， 获取指定的 condition
		var condition *myAppsv1.Condition
		for index := range md.Status.Conditions {
			// 4. 是否获取到
			if md.Status.Conditions[index].Type == conditionType {
				// 4.1. 获取到，
				condition = &md.Status.Conditions[index]
				// todo: break?
			}
		}

		if condition != nil {
			// 4.1.1. 获取当前线上的 condition 状态， 与存储的 condition 进行比较，如果相同，则跳过，不同则替换
			if condition.Status != status ||
				condition.Message != message ||
				condition.Reason != reason {
				condition.Status = status
				condition.Message = message
				condition.Reason = reason
			}
		} else {
			// 4.2. 没获取到， 创建这个 condition， 更新到 condition 中
			md.Status.Conditions = append(md.Status.Conditions,
				createCondition(conditionType, message, status, reason))
		}
	}
	// 5. 继续处理其他的 conditions

	message, reason, phase, success := isSuccess(md.Status.Conditions)
	if success {
		// 6.1. 如果所有的 conditions 的状态都为成功， 则更新总的 status 为成功。
		md.Status.Message = myAppsv1.StatusMessageSuccess
		md.Status.Reason = myAppsv1.StatusReasonSuccess
		md.Status.Phase = myAppsv1.StatusPhaseCompile
	} else {
		// 6.2. 变量所有的 conditions 状态， 如果有任意一个 condition 不是完成的状态， 则将这个状态更新到总的 status 中。等待一定时间再次入队。
		md.Status.Message = message
		md.Status.Reason = reason
		md.Status.Phase = phase
	}
	// 7. 执行更新
	return success, r.Client.Status().Update(ctx, md)
}

func isSuccess(conditions []myAppsv1.Condition) (message, reason, phase string, success bool) {
	if len(conditions) == 0 {
		return "", "", "", false
	}

	for _, condition := range conditions {
		if condition.Status == myAppsv1.ConditionStatusFalse {
			return condition.Message, condition.Reason, condition.Type, false
		}
	}

	return "", "", "", true
}

func createCondition(conditionType, message, status, reason string) myAppsv1.Condition {
	return myAppsv1.Condition{
		Type:               conditionType,
		Message:            message,
		Status:             status,
		Reason:             reason,
		LastTransitionTime: metav1.NewTime(time.Now()),
	}
}

// 需要是幂等的， 可以多次执行， 不管是否存在， 如果存在就删除， 不存在就什么也不做
// 只是删除对应conditions， 不做更多的操作
func (r *HxzDeploymentReconciler) deleteStatus(ctx context.Context, md *myAppsv1.HxzDeployment, conditionType string) {
	// 1. 遍历 conditions
	for index := range md.Status.Conditions {
		// 2. 找到要删除的对象
		if md.Status.Conditions[index].Type == conditionType {
			// 3.1. 执行删除
			md.Status.Conditions = deleteCondition(md.Status.Conditions, index)
			// todo: break?
		}
	}
}

func deleteCondition(conditions []myAppsv1.Condition, index int) []myAppsv1.Condition {
	// 切片中的元素顺序不敏感
	// 1. 要删除的元素的索引值不能大于切片长度
	if index > len(conditions) {
		return []myAppsv1.Condition{}
	}

	// 2. 如果切片长度为 1， 且索引值为0，直接清空
	if len(conditions) == 1 && index == 0 {
		return conditions[:0]
	}

	// 3. 如果长度-1 等于索引值， 删除最后一个元素
	if len(conditions)-1 == index {
		return conditions[:len(conditions)-1]
	}

	// 4. 交换索引位置的元素和最后一个元素， 删除最后一个元素
	conditions[index], conditions[len(conditions)-1] = conditions[len(conditions)-1], conditions[index]
	return conditions[:len(conditions)-1]

}
