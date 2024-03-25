package update

// import (
// 	"context"
// 	"github.com/huangxinzhong/huangxinzhong-deployment/test/framework"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// 	"k8s.io/apimachinery/pkg/runtime/schema"
// 	"k8s.io/client-go/dynamic"
// 	"k8s.io/client-go/kubernetes"
// 	"time"
// )

// // UpdateI2NHxzDeployment 测试从 ingress 模式更新为 nodeport 模式
// func UpdateI2NHxzDeployment(ctx *framework.TestContext, f *framework.Framework) {
// 	var (
// 		// 1. 准备测试数据
// 		ctFilePath       = "update/testdata/update-ingress.yaml"
// 		ctUpdateFilePath = "update/testdata/update-i2n.yaml"
// 		obj              = &unstructured.Unstructured{Object: make(map[string]interface{})}
// 		objUpdate        = &unstructured.Unstructured{Object: make(map[string]interface{})}
// 		dc               dynamic.Interface
// 		cs               *kubernetes.Clientset
// 		NameSpace        = "default"
// 		// 3. 准备测试用到的全局变量
// 		hxzGVR = schema.GroupVersionResource{
// 			Group:    "apps.huangxinzhong.com",
// 			Version:  "v1",
// 			Resource: "hxzdeployments",
// 		}
// 		err error
// 	)

// 	// 4. 初始化测试用到的全局变量
// 	BeforeEach(func() {
// 		// 2. 加载测试数据
// 		err = f.LoadYamlToUnstructured(ctFilePath, obj)
// 		Expect(err).Should(BeNil())

// 		err = f.LoadYamlToUnstructured(ctUpdateFilePath, objUpdate)
// 		Expect(err).Should(BeNil())

// 		dc = ctx.CreateDynamicClient()
// 		cs = ctx.CreateClientSet()
// 	})
// 	Context("Update hxzdeployment mod ingress to nodeport", func() {
// 		It("should be create mod ingress success", func() {
// 			_, err = dc.Resource(hxzGVR).Namespace(NameSpace).Create(context.TODO(), obj, metav1.CreateOptions{})
// 			Expect(err).Should(BeNil())

// 			By("Sleep 1 second wait creating done")
// 			time.Sleep(time.Second)
// 		})
// 		It("should be exist hxzdeployment", func() {
// 			_, err = dc.Resource(hxzGVR).Namespace(NameSpace).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
// 			Expect(err).Should(BeNil())
// 		})
// 		It("should be exist deployment", func() {
// 			_, err = cs.AppsV1().Deployments(NameSpace).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
// 			Expect(err).Should(BeNil())
// 		})
// 		It("should be exist service", func() {
// 			_, err = cs.CoreV1().Services(NameSpace).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
// 			Expect(err).Should(BeNil())
// 		})
// 		It("should be exist ingress", func() {
// 			_, err = cs.NetworkingV1().Ingresses(NameSpace).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
// 			Expect(err).Should(BeNil())
// 		})

// 		It("should be update to nodeport success", func() {
// 			var md *unstructured.Unstructured
// 			md, err = dc.Resource(hxzGVR).Namespace(NameSpace).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
// 			Expect(err).Should(BeNil())

// 			objUpdate.SetResourceVersion(md.GetResourceVersion())
// 			_, err = dc.Resource(hxzGVR).Namespace(NameSpace).Update(context.TODO(), objUpdate, metav1.UpdateOptions{})
// 			Expect(err).Should(BeNil())

// 			By("Sleep 1 second wait creating done")
// 			time.Sleep(time.Second)
// 		})
// 		It("should not be exist ingress", func() {
// 			_, err = cs.NetworkingV1().Ingresses(NameSpace).Get(context.TODO(), objUpdate.GetName(), metav1.GetOptions{})
// 			Expect(err).ShouldNot(BeNil())
// 		})
// 	})

// 	Context("Delete hxzdeployment i2n", func() {
// 		It("should be delete mod ingress success", func() {
// 			err = dc.Resource(hxzGVR).Namespace(NameSpace).Delete(context.TODO(), objUpdate.GetName(), metav1.DeleteOptions{})
// 			Expect(err).Should(BeNil())
// 			By("Sleep 1 second wait deleting done")
// 			time.Sleep(time.Second)
// 		})
// 		It("should not be exist hxzdeployment", func() {
// 			_, err = cs.AppsV1().Deployments(NameSpace).Get(context.TODO(), objUpdate.GetName(), metav1.GetOptions{})
// 			Expect(err).ShouldNot(BeNil())
// 		})
// 		It("should not be exist deployment", func() {
// 			_, err = cs.AppsV1().Deployments(NameSpace).Get(context.TODO(), objUpdate.GetName(), metav1.GetOptions{})
// 			Expect(err).ShouldNot(BeNil())
// 		})
// 		It("should not be exist service", func() {
// 			_, err = cs.CoreV1().Services(NameSpace).Get(context.TODO(), objUpdate.GetName(), metav1.GetOptions{})
// 			Expect(err).ShouldNot(BeNil())
// 		})
// 		It("should not be exist ingress", func() {
// 			_, err = cs.NetworkingV1().Ingresses(NameSpace).Get(context.TODO(), objUpdate.GetName(), metav1.GetOptions{})
// 			Expect(err).ShouldNot(BeNil())
// 		})
// 	})
// }
