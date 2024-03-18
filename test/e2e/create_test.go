//go:build e2e
// +build e2e

package e2e

import (
	"github.com/huangxinzhong/huangxinzhong-deployment/test/e2e/create"
)

var _ = fmw.Describe("Create huangxinzhong deployment mod ingress", create.CreateIngressHxzDeployment)
var _ = fmw.Describe("Create huangxinzhong deployment mod nodeport", create.CreateNodePortHxzDeployment)
