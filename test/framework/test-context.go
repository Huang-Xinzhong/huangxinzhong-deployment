package framework

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// 1. 定义一个测试入口函数 Describe，接受测试的模式以及 contextFunc
// 1.1 这里边回调用 context 创建方法来创建 context
// 1.2 这个 context 里面会有执行一些我们期望的行为
// 2. 这个contextFunc 的签名符合 func(ctx *TestContext, f *Framework)
// 3. 这个 contextFunc 的函数体就是测试函数的内容本身
// 4. 由于这个 contextFunc 的参数中有 ctx 入参， name 在执行测试函数体的时候， 就可以使用 ctx 种的内容或方法。

type TestContext struct {
	Name      string
	Namespace string
	Config    *rest.Config
	MasterIP  string
}

type ContextFunc func(ctx *TestContext, f *Framework)

// 如果不用动态的 client， name 我们访问这些资源的时候， 就需要:
// 1. 自己创建 rest api 的请求
// 2. 获取对应资源的 client sdk

// CreateDynamicClient 创建动态 client， 用来访问自定义或者后安装资源
func (tc *TestContext) CreateDynamicClient() dynamic.Interface {
	By("Create a Dynamic Client")
	c, err := dynamic.NewForConfig(tc.Config)
	if err != nil {
		Expect(err).Should(BeNil())
	}
	return c
}

// CreateClientSet 创建 clientset， 用来访问内置资源
func (tc *TestContext) CreateClientSet() *kubernetes.Clientset {
	By("Create a ClientSet client")
	c, err := kubernetes.NewForConfig(tc.Config)
	if err != nil {
		Expect(err).Should(BeNil())
	}
	return c
}
