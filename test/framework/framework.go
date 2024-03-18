package framework

import (
	"context"
	"flag"
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

// var DefaultStartTimeout = float64(60 * 60)
var DefaultStartTimeout = float64(60 * 60)
var nsRegex = regexp.MustCompile("[^a-z0-9]")

type Framework struct {
	Config        *Config
	ClusterConfig *ClusterConfig
	factory       Factory              // 工厂对象， 提供创建 provider 的方法， eq: factory := Factory{}
	provider      ClusterProvider      // 存储当前 framework 对象中实现的 provider
	client        kubernetes.Interface // 用来链接创建的 k8s 集群
	configFile    string               // 配置文件的路径
	initTimeout   float64              // 启动时候， 包括安装集群和依赖集本程序的超时时间
}

func NewFramework() *Framework {
	return &Framework{}
}

// Flags 解析命令行参数:  --config --timeout
func (f *Framework) Flags() *Framework {
	flag.StringVar(&f.configFile, "config", "config", "config file to used")
	flag.Float64Var(&f.initTimeout, "startup-timeout", DefaultStartTimeout, "startup timeout")
	flag.Parse()
	return f
}

// LoadConfig 加载配置文件到 fmw
func (f *Framework) LoadConfig(writer io.Writer) *Framework {
	// 1. 创建 config 对象
	config := NewConfig()

	// 2. 加载配置文件内容到 config 对象中
	if err := config.Load(f.configFile); err != nil {
		panic(err)
	}
	// 3. 将传入的 writer 应用到 config 中, 将 config 加入到 fmw 中
	return f.WithConfig(config.WithWriter(writer))
}

func (f *Framework) SynchronizedBeforeSuite(initFunc func()) *Framework {
	if initFunc == nil {
		initFunc = func() {
			// 1. 安装环境
			ginkgo.By("Deploying test environment")
			if err := f.DeployTestEnvironment(); err != nil {
				panic(err)
			}
			// 2. 初始化环境和访问的授权， 也就是创建 kubelet， 访问需要的 config
			ginkgo.By("kubectl switch context")
			// TODO: kubectl config
			kubectlConfig := NewKubectlConfig(f.Config)
			if err := kubectlConfig.SetContext(f.ClusterConfig); err != nil {
				panic(err)
			}

			// 退出前清理context
			defer func() {
				ginkgo.By("Kubectl reverting context")
				if !f.Config.Sub("cluster").Sub("kind").GetBool("retain") {
					_ = kubectlConfig.DeleteContext(f.ClusterConfig)
				}
			}()

			// 3. 安装依赖和我们的程序
			ginkgo.By("Preparing install steps")
			install := NewInstaller(f.Config)
			ginkgo.By("Executing install steps")
			if err := install.Install(); err != nil {
				panic(err)
			}

		}
	}
	ginkgo.SynchronizedBeforeSuite(
		func() []byte {
			initFunc()
			return nil
		},
		func(_ []byte) {},
		f.initTimeout,
	)
	return f
}

func (f *Framework) SynchronizedAfterSuite(destroyFunc func()) *Framework {
	if destroyFunc == nil {
		destroyFunc = func() {
			// 回收测试环境
			if err := f.DestroyTestEnvironment(); err != nil {
				panic(err)
			}
		}
	}
	ginkgo.SynchronizedAfterSuite(
		func() {},
		destroyFunc, f.initTimeout,
	)
	return f
}

func (f *Framework) MRun(m *testing.M) {
	rand.Seed(time.Now().UnixNano()) // 优化随机数
	os.Exit(m.Run())                 // 执行真正的 TestMain
}

func (f *Framework) Run(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	var r []ginkgo.Reporter
	r = append(r, reporters.NewJUnitReporter("e2e.xml"))
	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "e2e", r)
}

func (f *Framework) WithConfig(config *Config) *Framework {
	f.Config = config
	return f
}

//// DeployTestEnvironment 创建测试环境， 并且获取访问集群的配置及 client
//func (f *Framework) DeployTestEnvironment() error {
//	// 1. 检查 f.Config
//	if f.Config == nil {
//		return field.Invalid(
//			field.NewPath("Config"),
//			nil,
//			"not initial config object",
//		)
//	}
//
//	// 2. 创建 provider
//	ginkgo.By("Getting env provider")
//	var err error
//	if f.provider, err = f.factory.Provider(f.Config); err != nil {
//		return err
//	}
//
//	// 3. 执行 provider 实现的 validate 方法，验证 config
//	ginkgo.By("Validate config for provider")
//	if err := f.provider.Validate(f.Config); err != nil {
//		return err
//	}
//
//	// 4. 执行 provider 实现的 deploy 方法， 创建集群
//	ginkgo.By("Deploying test env")
//	clusterConfig, err := f.provider.Deploy(f.Config)
//	if err != nil {
//		return err
//	}
//	f.ClusterConfig = &clusterConfig
//
//	// 5. 创建 client， 用于执行测试用例的时候使用
//	if f.client, err = kubernetes.NewForConfig(f.ClusterConfig.Rest); err != nil {
//		return err
//	}
//	return nil
//}

// DeployTestEnvironment 创建测试环境，并且获取访问集群的配置及client
func (f *Framework) DeployTestEnvironment() error {
	// 1. 检查 f.config
	if f.Config == nil {
		return field.Invalid(
			field.NewPath("config"),
			nil,
			"Not inital config object")
	}
	// 2. 创建provider
	ginkgo.By("Getting env provider")
	var err error
	if f.provider, err = f.factory.Provider(f.Config); err != nil {
		return err
	}
	// 3. 执行 provider 实现的 validate 方法验证 config
	ginkgo.By("Validate config for provider")
	if err := f.provider.Validate(f.Config); err != nil {
		return err
	}

	// 4. 执行 provider 实现的 deploy 方法，创建集群
	ginkgo.By("Deploying test env")
	clusterConfig, err := f.provider.Deploy(f.Config)
	if err != nil {
		return err
	}
	f.ClusterConfig = &clusterConfig

	// 5. 创建 client，用于执行测试用例的时候使用
	if f.client, err = kubernetes.NewForConfig(f.ClusterConfig.Rest); err != nil {
		return err
	}

	return nil
}

// DestroyTestEnvironment 销毁测试环境， 此方法要在执行过 DeployTestEnvironment 方法之后执行
func (f *Framework) DestroyTestEnvironment() error {
	// 1. 检查 f.Config
	if f.Config == nil {
		return field.Invalid(
			field.NewPath("Config"),
			nil,
			"not initial config object",
		)
	}

	// 2. 检查 provider
	if f.provider == nil {
		return fmt.Errorf("f.provider is nil")
	}

	// 3. 执行 provider 的 destroy 方法来销毁环境
	ginkgo.By("Destroying test env")
	if err := f.provider.Destroy(f.Config); err != nil {
		return err
	}

	// 4. 清空 f.provider, 保护销毁函数呗多次执行而报错
	f.provider = nil

	return nil
}

// 加载测试文件内容到对象中
func (f *Framework) LoadYamlToUnstructured(ctFilePath string, obj *unstructured.Unstructured) error {
	data, err := os.ReadFile(ctFilePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &(obj.Object))
}

func (f *Framework) Describe(name string, ctxFunc ContextFunc) bool {
	// 整个函数， 实际上是调用 ginkgo 的 Describe
	return ginkgo.Describe(
		name,
		func() {
			// 1. 创建 testcontext
			ctx, err := f.createTestContext(name, false)
			if err != nil {
				ginkgo.Fail("Cannot create test context for " + name)
				return
			}

			// 2. 执行每次测试任务前， 来执行一些我们期望的动作， 如创建 namespace
			ginkgo.BeforeEach(func() {
				ctx2, err := f.createTestContext(name, true)
				if err != nil {
					ginkgo.Fail("Cannot create test context for " + name + " namespace " + ctx2.Namespace)
					return
				}
				ctx = ctx2
			})

			// 3. 执行每次测试任务之后要做的事情， 例如删除 testcontext
			ginkgo.AfterEach(func() {
				// 回收 testcontext
				_ = f.deleteTestContext(ctx)
			})
			// 4. 执行用户的测试函数
			ctxFunc(&ctx, f)
		},
	)
}

func (f *Framework) createTestContext(name string, nsCreate bool) (TestContext, error) {
	// 1. 创建 testcontext 对象

	tc := TestContext{}
	// 2. 检查 f 是否为空

	if f.Config == nil || f.ClusterConfig == nil {
		return tc, nil
		//return tc, field.Invalid(
		//	field.NewPath("config/clusterConfig"),
		//	nil,
		//	"not initial config object",
		//)
	}

	// 3. 填充字段
	tc.Name = name
	tc.Config = rest.CopyConfig(f.ClusterConfig.Rest)
	tc.MasterIP = f.ClusterConfig.MasterIP

	// 4. 判断参数，是否创建 namespace
	if nsCreate {
		// 4.1 如果创建， 使用 f.client 来创建 namespace
		// 4.1.1 处理 name， 将空格或下划线替换为'-'
		// 4.1.2 正则检查是否有其他非法字符
		// 4.1.3 自动生成 namespace 的机制
		prefix := nsRegex.ReplaceAllString(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ToLower(name),
					" ", "-"),
				"_", "-"),
			"")
		if len(prefix) > 30 {
			prefix = prefix[:30]
		}

		ns, err := f.client.CoreV1().Namespaces().Create(
			context.TODO(),
			&corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: prefix + "-",
				},
			},
			metav1.CreateOptions{},
		)
		if err != nil {
			return tc, err
		}
		tc.Namespace = ns.GetName()
	}
	// 5. 执行其他
	// 创建 sa/secret
	return tc, nil
}

func (f *Framework) deleteTestContext(ctx TestContext) error {
	// 删除创建的资源
	errs := field.ErrorList{}
	if ctx.Namespace != "" {
		if err := f.client.CoreV1().Namespaces().Delete(context.TODO(), ctx.Namespace, metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
			errs = append(errs, field.InternalError(field.NewPath("testcontext"), err))
		}
	}
	return errs.ToAggregate()
}
