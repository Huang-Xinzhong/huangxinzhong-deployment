//go:build e2e
// +build e2e

package e2e

import (
	"github.com/huangxinzhong/huangxinzhong-deployment/test/e2e/update"
)

var _ = fmw.Describe("Update huangxinzhong deployment mod ingress to nodeport", update.UpdateI2NHxzDeployment)
var _ = fmw.Describe("Update huangxinzhong deployment mod nodeport to ingress", update.UpdateN2IHxzDeployment)
