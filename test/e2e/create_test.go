//go:build e2e
// +build e2e

package e2e

import (
	"github.com/huangxinzhong/huangxinzhong-deployment/test/e2e/create"
)

// var _ = fmw.Describe("Create huangxinzhong deployment mod ingress", create.CreateIngressHxzDeployment)
var _ = fmw.Describe("Create huangxinzhong deployment mod ingress with tls", create.CreateIngressHxzDeploymentWithTls)

// var _ = fmw.Describe("Create huangxinzhong deployment mod nodeport", create.CreateNodeportHxzDeployment)
// var _ = fmw.Describe("Create hxzdeployment mod ingress default value", create.CreateIngressHxzDeploymentDefaultValue)
// var _ = fmw.Describe("Create hxzdeployment mod ingress must failed", create.CreateIngressHxzDeploymentMustFailed)
// var _ = fmw.Describe("Create hxzdeployment mod nodeport must failed", create.CreateNodeportHxzDeploymentMustFailed)
