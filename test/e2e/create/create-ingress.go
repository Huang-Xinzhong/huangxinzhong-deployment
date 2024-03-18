package create

import (
	"context"
	"github.com/huangxinzhong/huangxinzhong-deployment/test/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"time"
)

// CreateIngressHxzDeployment 测试创建 ingress 模式
func CreateIngressHxzDeployment(ctx *framework.TestContext, f *framework.Framework) {
	var (
		// 1. 准备测试数据
		ctFilePath = "create/testdata/create-ingress.yaml"
		obj        = &unstructured.Unstructured{Object: make(map[string]interface{})}
		dc         dynamic.Interface
		cs         *kubernetes.Clientset
		NameSpace  = "default"
		// 3. 准备测试用到的全局变量
		hxzGVR = schema.GroupVersionResource{
			Group:    "apps.huangxinzhong.com",
			Version:  "v1",
			Resource: "hxzdeployments",
		}
		err error
	)

	// 4. 初始化测试用到的全局变量
	BeforeEach(func() {
		// 2. 加载测试数据
		err = f.LoadYamlToUnstructured(ctFilePath, obj)
		Expect(err).Should(BeNil())
		dc = ctx.CreateDynamicClient()
		cs = ctx.CreateClientSet()
	})
	Context("Create hxzdeployment mod ingress", func() {
		It("should be create mod ingress success", func() {
			_, err = dc.Resource(hxzGVR).Namespace(NameSpace).Create(context.TODO(), obj, metav1.CreateOptions{})
			Expect(err).Should(BeNil())
			By("Sleep 1 second wait creating done")
			time.Sleep(time.Second)
		})
		It("should be exist hxzdeployment", func() {
			_, err = dc.Resource(hxzGVR).Namespace(NameSpace).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
			Expect(err).Should(BeNil())
		})
		It("should be exist deployment", func() {
			_, err = cs.AppsV1().Deployments(NameSpace).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
			Expect(err).Should(BeNil())
		})
		It("should be exist service", func() {
			_, err = cs.CoreV1().Services(NameSpace).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
			Expect(err).Should(BeNil())
		})
		It("should be exist ingress", func() {
			_, err = cs.NetworkingV1().Ingresses(NameSpace).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
			Expect(err).Should(BeNil())
		})
	})

	Context("Delete hxzdeployment mod ingress", func() {
		It("should be delete mod ingress success", func() {
			err = dc.Resource(hxzGVR).Namespace(NameSpace).Delete(context.TODO(), obj.GetName(), metav1.DeleteOptions{})
			Expect(err).Should(BeNil())
			By("Sleep 1 second wait deleting done")
			time.Sleep(time.Second)
		})
		It("should not be exist hxzdeployment", func() {
			_, err = cs.AppsV1().Deployments(NameSpace).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
			Expect(err).ShouldNot(BeNil())
		})
		It("should not be exist deployment", func() {
			_, err = cs.AppsV1().Deployments(NameSpace).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
			Expect(err).ShouldNot(BeNil())
		})
		It("should not be exist service", func() {
			_, err = cs.CoreV1().Services(NameSpace).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
			Expect(err).ShouldNot(BeNil())
		})
		It("should not be exist ingress", func() {
			_, err = cs.NetworkingV1().Ingresses(NameSpace).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
			Expect(err).ShouldNot(BeNil())
		})
	})
}
