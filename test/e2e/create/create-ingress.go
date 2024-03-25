package create

import (
	"context"
	"fmt"
	"time"

	// networkv1 "k8s.io/api/networking/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	// "k8s.io/client-go/kubernetes"

	"github.com/huangxinzhong/huangxinzhong-deployment/test/framework"
)

// 真正的测试函数

// 测试创建Ingress模式
// func CreateIngressHxzDeployment(ctx *framework.TestContext, f *framework.Framework) {
// 	var (
// 		// 1. 准备测试数据
// 		ctFilePath = "create/testdata/create-ingress.yaml"
// 		obj        = &unstructured.Unstructured{Object: make(map[string]interface{})}
// 		dc         dynamic.Interface
// 		cs         *kubernetes.Clientset

// 		// 3. 准备测试用到的全局变量
// 		HxzGVR = schema.GroupVersionResource{
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
// 	Context("Create Hxzdeployment mod ingress", func() {
// 		It("Should be create mod ingress success", func() {
// 			_, err = dc.Resource(HxzGVR).Namespace("default").Create(context.TODO(), obj, metav1.CreateOptions{})
// 			Expect(err).Should(BeNil())

// 			By("Sleep 1 second wait creating done")
// 			time.Sleep(time.Second)
// 		})
// 		It("Should be exist Hxzdeployment", func() {
// 			_, err = dc.Resource(HxzGVR).Namespace("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
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
// 		It("Should be exist ingress", func() {
// 			_, err = cs.NetworkingV1().Ingresses("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
// 			Expect(err).Should(BeNil())
// 		})
// 	})

// 	Context("Delete Hxzdeployment mod ingress", func() {
// 		It("Should be delete mod ingress success", func() {
// 			err = dc.Resource(HxzGVR).Namespace("default").Delete(context.TODO(), obj.GetName(), metav1.DeleteOptions{})
// 			Expect(err).Should(BeNil())

// 			By("Sleep 3 second wait deleting done")
// 			time.Sleep(3 * time.Second)
// 		})
// 		It("Should not be exist Hxzdeployment", func() {
// 			_, err = dc.Resource(HxzGVR).Namespace("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
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

// 测试HTTPS
func CreateIngressHxzDeploymentWithTls(ctx *framework.TestContext, f *framework.Framework) {
	var (
		// 1. 准备测试数据
		ctFilePath = "create/testdata/create-ingress-with-tls.yaml"
		obj        = &unstructured.Unstructured{Object: make(map[string]interface{})}
		dc         dynamic.Interface
		// cs         *kubernetes.Clientset

		// 3. 准备测试用到的全局变量
		HxzGVR = schema.GroupVersionResource{
			Group:    "apps.huangxinzhong.com",
			Version:  "v1",
			Resource: "hxzdeployments",
		}
		// issuer
		// issuerGVR = schema.GroupVersionResource{
		// 	Group:    "cert-manager.io",
		// 	Version:  "v1",
		// 	Resource: "issuers",
		// }
		// certificate
		certGVR = schema.GroupVersionResource{
			Group:    "cert-manager.io",
			Version:  "v1",
			Resource: "certificates",
		}

		err error
	)
	BeforeEach(func() {
		// 2. 加载测试数据
		err = f.LoadYamlToUnstructured(ctFilePath, obj)
		Expect(err).Should(BeNil())

		// 4. 初始化测试用到的全局变量
		dc = ctx.CreateDynamicClient()
		// cs = ctx.CreateClientSet()
	})
	Context("Create Hxzdeployment mod ingress with tls", func() {
		It("Should be create mod ingress with tls success", func() {
			_, err := dc.Resource(HxzGVR).Namespace("default").Create(context.TODO(), obj, metav1.CreateOptions{})
			Expect(err).Should(BeNil())
			// fmt.Printf("debug >>>>>>>>> %#v\n", md)
			By("Sleep 3 second wait creating done")
			time.Sleep(3 * time.Second)
		})
		// It("Should be exist Hxzdeployment", func() {
		// 	_, err = dc.Resource(HxzGVR).Namespace("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
		// 	Expect(err).Should(BeNil())
		// })
		// It("Should be exist ingress, and have a tls setting", func() {
		// 	var ig *networkv1.Ingress
		// 	ig, err = cs.NetworkingV1().Ingresses("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
		// 	Expect(err).Should(BeNil())
		// 	Expect(len(ig.Spec.TLS)).To(Equal(1))
		// })
		// It("Should be exist issuer", func() {
		// 	md, err := dc.Resource(issuerGVR).Namespace("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
		// 	fmt.Printf("debug issuerResource >>>>>>>>>>> %#v\n", md)
		// 	fmt.Println("debug >>>>>>>>>>>>>>>>>>> obj.GetName() = ", obj.GetName())
		// 	Expect(err).Should(BeNil())
		// })
		It("Should be exist certificate", func() {
			md, err := dc.Resource(certGVR).Namespace("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
			fmt.Printf("debug issuerResource >>>>>>>>>>> %#v\n", md)
			fmt.Println("debug >>>>>>>>>>>>>>>>>>> obj.GetName() = ", obj.GetName())
			fmt.Println("debug >>>>>>>>>>>>>>>>>>> dc.Resource(certGVR).Namespace(\"default\").Get(context.TODO(), obj.GetName(), metav1.GetOptions{}) err = ", err)
			Expect(err).Should(BeNil())
		})
	})

	// Context("Delete Hxzdeployment mod ingress with tls", func() {
	// 	It("Should be delete mod ingress success", func() {
	// 		err = dc.Resource(HxzGVR).Namespace("default").Delete(context.TODO(), obj.GetName(), metav1.DeleteOptions{})
	// 		Expect(err).Should(BeNil())

	// 		By("Sleep 3 second wait deleting done")
	// 		time.Sleep(3 * time.Second)
	// 	})
	// 	It("Should not be exist Hxzdeployment", func() {
	// 		_, err = dc.Resource(HxzGVR).Namespace("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
	// 		Expect(err).ShouldNot(BeNil())
	// 	})
	// 	It("Should not be exist deployment", func() {
	// 		_, err = cs.AppsV1().Deployments("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
	// 		Expect(err).ShouldNot(BeNil())
	// 	})
	// 	It("Should not be exist service", func() {
	// 		_, err = cs.CoreV1().Services("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
	// 		Expect(err).ShouldNot(BeNil())
	// 	})
	// 	It("Should not be exist ingress", func() {
	// 		_, err = cs.NetworkingV1().Ingresses("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
	// 		Expect(err).ShouldNot(BeNil())
	// 	})
	// 	It("Should not be exist issuer", func() {
	// 		_, err = dc.Resource(issuerGVR).Namespace("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
	// 		Expect(err).ShouldNot(BeNil())
	// 	})
	// 	It("Should not be exist certificate", func() {
	// 		_, err = dc.Resource(certGVR).Namespace("default").Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
	// 		Expect(err).ShouldNot(BeNil())
	// 	})
	// })
}

// // 测试默认值设置
// func CreateIngressHxzDeploymentDefaultValue(ctx *framework.TestContext, f *framework.Framework) {
// 	var (
// 		// 1. 准备测试数据
// 		ctFileRepPath = "create/testdata/create-ingress-default-no-replicas.yaml"
// 		ctFileSvcPath = "create/testdata/create-ingress-default-no-serviceport.yaml"
// 		objRep        = &unstructured.Unstructured{Object: make(map[string]interface{})}
// 		objSvc        = &unstructured.Unstructured{Object: make(map[string]interface{})}
// 		dc            dynamic.Interface

// 		// 3. 准备测试用到的全局变量
// 		HxzGVR = schema.GroupVersionResource{
// 			Group:    "apps.huangxinzhong.com",
// 			Version:  "v1",
// 			Resource: "hxzdeployments",
// 		}
// 		err error
// 	)
// 	BeforeEach(func() {
// 		// 2. 加载测试数据
// 		err = f.LoadYamlToUnstructured(ctFileRepPath, objRep)
// 		Expect(err).Should(BeNil())

// 		err = f.LoadYamlToUnstructured(ctFileSvcPath, objSvc)
// 		Expect(err).Should(BeNil())

// 		// 4. 初始化测试用到的全局变量
// 		dc = ctx.CreateDynamicClient()
// 	})

// 	Context("Create Hxzdeployment mod ingress, but no replicas", func() {
// 		It("Create Hxzdeployment mod ingress success", func() {
// 			_, err = dc.Resource(HxzGVR).Namespace("default").Create(context.TODO(), objRep, metav1.CreateOptions{})
// 			Expect(err).Should(BeNil())

// 			By("Sleep 1 second wait creating done")
// 			time.Sleep(time.Second)
// 		})
// 		It("Should be exist Hxzdeployment, and have a default replicas", func() {
// 			var md *unstructured.Unstructured
// 			md, err = dc.Resource(HxzGVR).Namespace("default").Get(context.TODO(), objRep.GetName(), metav1.GetOptions{})
// 			Expect(err).Should(BeNil())

// 			data := md.UnstructuredContent()
// 			port, ok := data["spec"].(map[string]interface{})["replicas"].(int64)
// 			Expect(ok).To(Equal(true))
// 			Expect(port).To(Equal(int64(1)))
// 		})
// 	})

// 	Context("Create Hxzdeployment mod ingress, but no svcport", func() {
// 		It("Create Hxzdeployment mod ingress success", func() {
// 			_, err = dc.Resource(HxzGVR).Namespace("default").Create(context.TODO(), objSvc, metav1.CreateOptions{})
// 			Expect(err).Should(BeNil())

// 			By("Sleep 1 second wait creating done")
// 			time.Sleep(time.Second)
// 		})
// 		It("Should be exist Hxzdeployment, and have a default svcport", func() {
// 			var md *unstructured.Unstructured
// 			md, err = dc.Resource(HxzGVR).Namespace("default").Get(context.TODO(), objSvc.GetName(), metav1.GetOptions{})
// 			Expect(err).Should(BeNil())

// 			data := md.UnstructuredContent()
// 			port, ok := data["spec"].(map[string]interface{})["port"].(int64)
// 			Expect(ok).To(Equal(true))
// 			svcPort, svcOk := data["spec"].(map[string]interface{})["expose"].(map[string]interface{})["servicePort"].(int64)
// 			Expect(svcOk).To(Equal(true))
// 			Expect(port).To(Equal(svcPort))
// 		})
// 	})
// }

// func CreateIngressHxzDeploymentMustFailed(ctx *framework.TestContext, f *framework.Framework) {
// 	var (
// 		// 1. 准备测试数据
// 		ctFilePath = "create/testdata/create-ingress-error-no-domain.yaml"
// 		obj        = &unstructured.Unstructured{Object: make(map[string]interface{})}
// 		dc         dynamic.Interface

// 		// 3. 准备测试用到的全局变量
// 		HxzGVR = schema.GroupVersionResource{
// 			Group:    "apps.huangxinzhong.com",
// 			Version:  "v1",
// 			Resource: "Hxzdeployments",
// 		}
// 		err error
// 	)
// 	BeforeEach(func() {
// 		// 2. 加载测试数据
// 		err = f.LoadYamlToUnstructured(ctFilePath, obj)
// 		Expect(err).Should(BeNil())

// 		// 4. 初始化测试用到的全局变量
// 		dc = ctx.CreateDynamicClient()
// 	})
// 	Context("Create Hxzdeployment mod ingress no domain", func() {
// 		It("Should be create mod ingress no domain failed", func() {
// 			_, err = dc.Resource(HxzGVR).Namespace("default").Create(context.TODO(), obj, metav1.CreateOptions{})
// 			Expect(err).ShouldNot(BeNil())
// 		})
// 	})
// }

// // 测试默认值设置
// func CreateIngressHxzDeploymentDefaultValue(ctx *framework.TestContext, f *framework.Framework) {
// 	var (
// 		// 1. 准备测试数据
// 		ctFileRepPath = "create/testdata/create-ingress-default-no-replicas.yaml"
// 		ctFileSvcPath = "create/testdata/create-ingress-default-no-serviceport.yaml"
// 		objRep        = &unstructured.Unstructured{Object: make(map[string]interface{})}
// 		objSvc        = &unstructured.Unstructured{Object: make(map[string]interface{})}
// 		dc            dynamic.Interface

// 		// 3. 准备测试用到的全局变量
// 		HxzGVR = schema.GroupVersionResource{
// 			Group:    "apps.huangxinzhong.com",
// 			Version:  "v1",
// 			Resource: "Hxzdeployments",
// 		}
// 		err error
// 	)
// 	BeforeEach(func() {
// 		// 2. 加载测试数据
// 		err = f.LoadYamlToUnstructured(ctFileRepPath, objRep)
// 		Expect(err).Should(BeNil())

// 		err = f.LoadYamlToUnstructured(ctFileSvcPath, objSvc)
// 		Expect(err).Should(BeNil())

// 		// 4. 初始化测试用到的全局变量
// 		dc = ctx.CreateDynamicClient()
// 	})

// 	Context("Create Hxzdeployment mod ingress, but no replicas", func() {
// 		It("Create Hxzdeployment mod ingress success", func() {
// 			_, err = dc.Resource(HxzGVR).Namespace("default").Create(context.TODO(), objRep, metav1.CreateOptions{})
// 			Expect(err).Should(BeNil())

// 			By("Sleep 1 second wait creating done")
// 			time.Sleep(time.Second)
// 		})
// 		It("Should be exist Hxzdeployment, and have a default replicas", func() {
// 			var md *unstructured.Unstructured
// 			md, err = dc.Resource(HxzGVR).Namespace("default").Get(context.TODO(), objRep.GetName(), metav1.GetOptions{})
// 			Expect(err).Should(BeNil())

// 			data := md.UnstructuredContent()
// 			port, ok := data["spec"].(map[string]interface{})["replicas"].(int64)
// 			Expect(ok).To(Equal(true))
// 			Expect(port).To(Equal(int64(1)))
// 		})
// 	})

// 	Context("Create Hxzdeployment mod ingress, but no svcport", func() {
// 		It("Create Hxzdeployment mod ingress success", func() {
// 			_, err = dc.Resource(HxzGVR).Namespace("default").Create(context.TODO(), objSvc, metav1.CreateOptions{})
// 			Expect(err).Should(BeNil())

// 			By("Sleep 1 second wait creating done")
// 			time.Sleep(time.Second)
// 		})
// 		It("Should be exist Hxzdeployment, and have a default svcport", func() {
// 			var md *unstructured.Unstructured
// 			md, err = dc.Resource(HxzGVR).Namespace("default").Get(context.TODO(), objSvc.GetName(), metav1.GetOptions{})
// 			Expect(err).Should(BeNil())

// 			data := md.UnstructuredContent()
// 			port, ok := data["spec"].(map[string]interface{})["port"].(int64)
// 			Expect(ok).To(Equal(true))
// 			svcPort, svcOk := data["spec"].(map[string]interface{})["expose"].(map[string]interface{})["servicePort"].(int64)
// 			Expect(svcOk).To(Equal(true))
// 			Expect(port).To(Equal(svcPort))
// 		})
// 	})
// }

// func CreateIngressHxzDeploymentMustFailed(ctx *framework.TestContext, f *framework.Framework) {
// 	var (
// 		// 1. 准备测试数据
// 		ctFilePath = "create/testdata/create-ingress-error-no-domain.yaml"
// 		obj        = &unstructured.Unstructured{Object: make(map[string]interface{})}
// 		dc         dynamic.Interface

// 		// 3. 准备测试用到的全局变量
// 		HxzGVR = schema.GroupVersionResource{
// 			Group:    "apps.huangxinzhong.com",
// 			Version:  "v1",
// 			Resource: "Hxzdeployments",
// 		}
// 		err error
// 	)
// 	BeforeEach(func() {
// 		// 2. 加载测试数据
// 		err = f.LoadYamlToUnstructured(ctFilePath, obj)
// 		Expect(err).Should(BeNil())

// 		// 4. 初始化测试用到的全局变量
// 		dc = ctx.CreateDynamicClient()
// 	})
// 	Context("Create Hxzdeployment mod ingress no domain", func() {
// 		It("Should be create mod ingress no domain failed", func() {
// 			_, err = dc.Resource(HxzGVR).Namespace("default").Create(context.TODO(), obj, metav1.CreateOptions{})
// 			Expect(err).ShouldNot(BeNil())
// 		})
// 	})
// }
