//go:build e2e
// +build e2e

package e2e

import (
	"github.com/huangxinzhong/huangxinzhong-deployment/test/framework"
	"github.com/onsi/ginkgo"
	"testing"
)

var fmw = framework.NewFramework()

// 执行 go test 的时候， 会被先执行的内容
func TestMain(m *testing.M) {
	fmw.Flags(). // 解析命令行
			LoadConfig(ginkgo.GinkgoWriter). // 加载配置
			SynchronizedBeforeSuite(nil).    // 同步的， 在执行测试任务之前执行的内容
			SynchronizedAfterSuite(nil).     // 同步的， 在执行测试任务之后执行的内容
			MRun(m)
}

// 执行 go test 时候， 会被后执行， 也就是正常的测试用例
func TestE2E(t *testing.T) {
	fmw.Run(t)
}
