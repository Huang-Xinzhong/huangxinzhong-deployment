package framework

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/rest"
)

type ClusterConfig struct {
	Name     string       // 存储 kind create cluster name
	Rest     *rest.Config `json:"-"` // 链接创建的 k8s 的 client。 这个 client 比较低级。
	MasterIP string       // 集群的 master ip。 方便在一些需要直接和集群通讯的测试时使用
}

type ClusterProvider interface {
	Validate(config *Config) error
	Deploy(config *Config) (ClusterConfig, error)
	Destroy(config *Config) error
}

// 1. 定义工厂对象
type Factory struct {
}

// Provider 2. 工厂对象中创建不同实现的对象
func (f Factory) Provider(config *Config) (ClusterProvider, error) {
	var clusterProvider ClusterProvider

	// 1. 检查配置
	if config.Viper == nil {
		return clusterProvider, field.Invalid(
			field.NewPath("config"),
			nil,
			"not initial config object",
		)
	}
	// 2. 检查创建集群相关的 config
	if config.Sub("cluster") == nil {
		return clusterProvider, field.Invalid(
			field.NewPath("cluster"),
			nil,
			"not initial config object",
		)
	}
	cluster := config.Sub("cluster")

	// 3. 判断创建 k8s 集群的插件创建对象
	switch {
	case cluster.Sub("kind") != nil:
		kind := new(KindProvider)
		return kind, nil
	default:
		return clusterProvider, fmt.Errorf("not support provider: %#v", cluster.AllSettings())
	}

}

// 3.
