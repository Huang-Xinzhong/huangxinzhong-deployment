package create

// import (
// 	"context"
// 	"time"

// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// 	"k8s.io/apimachinery/pkg/runtime/schema"
// 	"k8s.io/client-go/dynamic"
// 	"k8s.io/client-go/kubernetes"

// 	"github.com/huangxinzhong/huangxinzhong-deployment/test/framework"
// )

// // 真正的测试函数

// // 测试创建Nodeport模式
// func CreateNodeportHxzDeployment(ctx *framework.TestContext, f *framework.Framework) {
// 	var (
// 		// 1. 准备测试数据
// 		ctFilePath = "create/testdata/create-nodeport.yaml"
// 		obj        = &unstructured.Unstructured{Object: make(map[string]interface{})}
// 		dc         dynamic.Interface
// 		cs         *kubernetes.Clientset

// 		// 3. 准备测试用到的全局变量
// 		hxzGVR = schema.GroupVersionResource{
// 			Group:    "apps.huangxinzhong.com",
// 			Version:  "v1",
// 			Resource: "hxzdeployments",
// 		}
// 		err error
// 	)
// 	BeforeEach(func() {
// 		// 2. 加载测试数据
// 		err = f.LoadYamlToUnstructured(ctFilePath, obj)
// 		Expect(err).Should(BeNil())

// 		// 4. 初始化测试用到的全局变量
// 		dc = ctx.CreateDynamicClient()
// 		cs = ctx.CreateClientSet()
// 	})
// 	Context("Create hxzdeployment mod nodeport", func() {
// 		It("Should be create mod nodeport success", func() {
// 			_, err = dc.Resource(hxzGVR).Namespace("default").Create(context.TODO(), obj, metav1.CreateOptions{})
// 			Expect(err).Should(BeNil())

// 			By("Sleep 1 second wait creating done")
// 			time.Sleep(time.Second)
// 		})
// 		It("Should be exist hxzdeployment", func() {
// 			_, err = dc.Resource(hxzGVR).Namespace("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
// 			Expect(err).Should(BeNil())
// 		})
// 		It("Should be exist deployment", func() {
// 			_, err = cs.AppsV1().Deployments("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
// 			Expect(err).Should(BeNil())
// 		})
// 		It("Should be exist service", func() {
// 			_, err = cs.CoreV1().Services("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
// 			Expect(err).Should(BeNil())
// 		})
// 		It("Should not be exist ingress", func() {
// 			_, err = cs.NetworkingV1().Ingresses("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
// 			Expect(err).ShouldNot(BeNil())
// 		})
// 	})

// 	Context("Delete hxzdeployment mod nodeport", func() {
// 		It("Should be delete mod nodeport success", func() {
// 			err = dc.Resource(hxzGVR).Namespace("default").Delete(context.TODO(), obj.GetName(), metav1.DeleteOptions{})
// 			Expect(err).Should(BeNil())

// 			By("Sleep 3 second wait deleting done")
// 			time.Sleep(3 * time.Second)
// 		})
// 		It("Should not be exist hxzdeployment", func() {
// 			_, err = dc.Resource(hxzGVR).Namespace("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
// 			Expect(err).ShouldNot(BeNil())
// 		})
// 		It("Should not be exist deployment", func() {
// 			_, err = cs.AppsV1().Deployments("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
// 			Expect(err).ShouldNot(BeNil())
// 		})
// 		It("Should not be exist service", func() {
// 			_, err = cs.CoreV1().Services("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
// 			Expect(err).ShouldNot(BeNil())
// 		})
// 		It("Should not be exist ingress", func() {
// 			_, err = cs.NetworkingV1().Ingresses("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
// 			Expect(err).ShouldNot(BeNil())
// 		})
// 	})
// }

// // func CreateNodeportHxzDeploymentMustFailed(ctx *framework.TestContext, f *framework.Framework) {
// // 	var (
// // 		// 1. 准备测试数据
// // 		ctFileNoPath = "create/testdata/create-nodeport-error-no-nodeport.yaml"
// // 		ctFileLiPath = "create/testdata/create-nodeport-error-li-30000.yaml"
// // 		ctFileGtPath = "create/testdata/create-nodeport-error-gt-32767.yaml"
// // 		objNo        = &unstructured.Unstructured{Object: make(map[string]interface{})}
// // 		objLi        = &unstructured.Unstructured{Object: make(map[string]interface{})}
// // 		objGt        = &unstructured.Unstructured{Object: make(map[string]interface{})}
// // 		dc           dynamic.Interface

// // 		// 3. 准备测试用到的全局变量
// // 		hxzGVR = schema.GroupVersionResource{
// // 			Group:    "apps.huangxinzhong.com",
// // 			Version:  "v1",
// // 			Resource: "hxzdeployments",
// // 		}
// // 		err error
// // 	)
// // 	BeforeEach(func() {
// // 		// 2. 加载测试数据
// // 		err = f.LoadYamlToUnstructured(ctFileNoPath, objNo)
// // 		Expect(err).Should(BeNil())
// // 		err = f.LoadYamlToUnstructured(ctFileLiPath, objLi)
// // 		Expect(err).Should(BeNil())
// // 		err = f.LoadYamlToUnstructured(ctFileGtPath, objGt)
// // 		Expect(err).Should(BeNil())

// // 		// 4. 初始化测试用到的全局变量
// // 		dc = ctx.CreateDynamicClient()
// // 	})
// // 	Context("Create hxzdeployment mod nodeport", func() {
// // 		It("Should be create mod nodeport no nodeport, must failed", func() {
// // 			_, err = dc.Resource(hxzGVR).Namespace("default").Create(context.TODO(), objNo, metav1.CreateOptions{})
// // 			Expect(err).ShouldNot(BeNil())
// // 		})
// // 		It("Should be create mod nodeport li nodeport, must failed", func() {
// // 			_, err = dc.Resource(hxzGVR).Namespace("default").Create(context.TODO(), objLi, metav1.CreateOptions{})
// // 			Expect(err).ShouldNot(BeNil())
// // 		})
// // 		It("Should be create mod nodeport gt nodeport, must failed", func() {
// // 			_, err = dc.Resource(hxzGVR).Namespace("default").Create(context.TODO(), objGt, metav1.CreateOptions{})
// // 			Expect(err).ShouldNot(BeNil())
// // 		})
// // 	})
// // }

// func CreateNodeportHxzDeploymentMustFailed(ctx *framework.TestContext, f *framework.Framework) {
// 	var (
// 		// 1. 准备测试数据
// 		ctFileNoPath = "create/testdata/create-nodeport-error-no-nodeport.yaml"
// 		ctFileLiPath = "create/testdata/create-nodeport-error-li-30000.yaml"
// 		ctFileGtPath = "create/testdata/create-nodeport-error-gt-32767.yaml"
// 		objNo        = &unstructured.Unstructured{Object: make(map[string]interface{})}
// 		objLi        = &unstructured.Unstructured{Object: make(map[string]interface{})}
// 		objGt        = &unstructured.Unstructured{Object: make(map[string]interface{})}
// 		dc           dynamic.Interface

// 		// 3. 准备测试用到的全局变量
// 		hxzGVR = schema.GroupVersionResource{
// 			Group:    "apps.huangxinzhong.com",
// 			Version:  "v1",
// 			Resource: "hxzdeployments",
// 		}
// 		err error
// 	)
// 	BeforeEach(func() {
// 		// 2. 加载测试数据
// 		err = f.LoadYamlToUnstructured(ctFileNoPath, objNo)
// 		Expect(err).Should(BeNil())
// 		err = f.LoadYamlToUnstructured(ctFileLiPath, objLi)
// 		Expect(err).Should(BeNil())
// 		err = f.LoadYamlToUnstructured(ctFileGtPath, objGt)
// 		Expect(err).Should(BeNil())

// 		// 4. 初始化测试用到的全局变量
// 		dc = ctx.CreateDynamicClient()
// 	})
// 	Context("Create hxzdeployment mod nodeport", func() {
// 		It("Should be create mod nodeport no nodeport, must failed", func() {
// 			_, err = dc.Resource(hxzGVR).Namespace("default").Create(context.TODO(), objNo, metav1.CreateOptions{})
// 			Expect(err).ShouldNot(BeNil())
// 		})
// 		It("Should be create mod nodeport li nodeport, must failed", func() {
// 			_, err = dc.Resource(hxzGVR).Namespace("default").Create(context.TODO(), objLi, metav1.CreateOptions{})
// 			Expect(err).ShouldNot(BeNil())
// 		})
// 		It("Should be create mod nodeport gt nodeport, must failed", func() {
// 			_, err = dc.Resource(hxzGVR).Namespace("default").Create(context.TODO(), objGt, metav1.CreateOptions{})
// 			Expect(err).ShouldNot(BeNil())
// 		})
// 	})
// }
