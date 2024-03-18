package framework

import (
	"bytes"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/exec"
)

type KindConfig struct {
	Name   string `json:"name"`   // 创建 cluster 的参数使用
	Config string `json:"config"` // Config 会当做 kind create cluster --config 中传入
	Retain bool   `json:"retain"` // 在执行完测试任务后是否保留

}

type KindProvider struct{}

// 检查这个对象是否实现了接口 ClusterProvider, 这种方法也推荐在一切实现某个接口对象来检查
var _ ClusterProvider = &KindProvider{}

func (k *KindProvider) Validate(config *Config) error {
	// 1 获取配置
	if config == nil {
		return field.Invalid(field.NewPath("config"), nil, "not initial config object")
	}

	kindConfig := config.Sub("cluster").Sub("kind")
	root := field.NewPath("cluster", "kind")
	if kindConfig == nil {
		return field.Invalid(root, nil, "Config does not have kind configuration")
	}

	// 2 检查必要项
	if kindConfig.GetString("name") == "" {
		// 3 设置默认
		kindConfig.Set("name", "e2e")
	}

	return nil
}

func (k *KindProvider) Deploy(config *Config) (ClusterConfig, error) {
	clusterConfig := ClusterConfig{}
	// 1 获取配置
	kindConfig, err := getKindConfig(config.Sub("cluster").Sub("kind"))
	if err != nil {
		return clusterConfig, err
	}

	var kubeConfigFile string
	// 2 确认是否存在 cluster
	if kindConfig.Name != "" {
		// 判断集群是否存在的命令: kind get kubeconfig --name <kindConfig.Name>
		output := &bytes.Buffer{} // 用来接受配置文件的内容
		cmd := exec.Command(
			"kind",
			"get",
			"kubeconfig",
			"--name",
			kindConfig.Name,
		)
		cmd.Stdout = output
		cmd.Stderr = config.Stderr

		if err := cmd.Run(); err == nil {
			// 2.1 存在：生成访问 k8s 集群的 config 文件
			if err := os.WriteFile(KubeconfigTempFile, output.Bytes(), os.ModePerm); err != nil {
				return clusterConfig, err
			}
			kubeConfigFile = KubeconfigTempFile
		}
	}

	if kubeConfigFile == "" {
		// 2.2 不存在: 创建， 并返回访问 k8s 集群的配置文件
		// 集群命令: kind create cluster --kubeconfig <KubeconfigTempFile> --config <KindConfigTempFile>
		subCommand := []string{"create", "cluster", "--kubeconfig", KubeconfigTempFile}
		if kindConfig.Config != "" {
			if err := os.WriteFile(
				KindConfigTempFile,
				[]byte(kindConfig.Config),
				os.ModePerm,
			); err != nil {
				return clusterConfig, err
			}
			defer func() { _ = os.Remove(KindConfigTempFile) }()
			subCommand = append(subCommand, "--config", KindConfigTempFile)
		}
		subCommand = append(subCommand, "--name", kindConfig.Name)
		cmd := exec.Command("kind", subCommand...)
		cmd.Stdout = config.Stdout
		cmd.Stderr = config.Stderr
		if err := cmd.Run(); err != nil {
			return clusterConfig, err
		}
		kubeConfigFile = KubeconfigTempFile
	}
	defer func() { _ = os.Remove(kubeConfigFile) }() // 退出函数之前清空 kubeconfig 文件

	// 3 创建 Cluster config
	clusterConfig.Name = kindConfig.Name
	if clusterConfig.Rest, err = clientcmd.BuildConfigFromFlags("", kubeConfigFile); err != nil {
		return clusterConfig, err
	}

	clusterConfig.MasterIP = clusterConfig.Rest.Host
	//ginkgo.By(fmt.Sprintf("clusterConfig: %#v", clusterConfig))
	return clusterConfig, nil
}

func (k *KindProvider) Destroy(config *Config) error {
	// 1 获取配置
	kindConfig, err := getKindConfig(config.Sub("cluster").Sub("kind"))
	if err != nil {
		return err
	}
	if kindConfig.Retain {
		// 2 保留就退出
		return nil
	}

	// 3 不保留就销毁
	// 销毁集群的命令 kind delete cluster --name <kindConfig.Name>
	cmd := exec.Command(
		"kind",
		"delete",
		"cluster",
		"--name",
		kindConfig.Name,
	)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func getKindConfig(config *viper.Viper) (KindConfig, error) {
	kindConfig := KindConfig{}

	if config == nil {
		return kindConfig, field.Invalid(
			field.NewPath("cluster", "kind"),
			nil,
			"not initial config object",
		)
	}

	if err := config.Unmarshal(&kindConfig); err != nil {
		return kindConfig, err
	}

	if kindConfig.Name == "" {
		kindConfig.Name = "e2e"
	}
	return kindConfig, nil
}
