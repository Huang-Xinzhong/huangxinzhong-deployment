package framework

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"path/filepath"
)

var ConfigFileGroupResource = schema.GroupResource{
	Group:    "",
	Resource: "config",
}

type Config struct {
	*viper.Viper // 动态处理配置文件的工具

	Stdout io.Writer
	Stderr io.Writer
}

func NewConfig() *Config {
	return &Config{
		Viper:  viper.New(),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Load 从文件加载配置内容到 config 对象中
func (c *Config) Load(fileName string) error {
	// 1 设置文件名
	c.SetConfigName(filepath.Base(fileName))
	c.AddConfigPath(filepath.Dir(fileName))

	if err := c.ReadInConfig(); err != nil {
		ext := filepath.Ext(fileName)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok && ext != "" {
			c.SetConfigName(filepath.Base(fileName[:len(fileName)-len(ext)]))
			err = c.ReadInConfig()
		}
		if err != nil {
			switch err.(type) {
			case viper.ConfigFileNotFoundError:
				return errors.NewNotFound(ConfigFileGroupResource, fmt.Sprintf("config file \"%s\" not found", fileName))
			case viper.UnsupportedConfigError:
				return errors.NewBadRequest("not using a supported file format")
			default:
				return err
			}
		}
	}
	return nil
}

func (c *Config) WithWriter(std io.Writer) *Config {
	c.Stdout = std
	c.Stderr = std
	return c
}
