package framework

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type KubectlConfig struct {
	previousContext string
	config          *Config // 主要是使用里面的 stdout 和 stderr
}

var contextNameRegex = regexp.MustCompile("[^a-zA-Z0-9_]")

func NewKubectlConfig(config *Config) *KubectlConfig {
	return &KubectlConfig{
		config: config,
	}
}

func (k *KubectlConfig) SetContext(config *ClusterConfig) error {
	contextName := k.getContextName(config)
	// 1. 获取当前的 context 保存起来(如果有)
	// 获取命令: kubectl config current-context
	cmd := k.Command("kubectl", "config", "current-context")
	currentContext := &bytes.Buffer{}
	cmd.Stdout = currentContext
	err := cmd.Run()
	// 如果执行成功， 说明存在 currentContext，就把他保存起来，在DeleteContext 的时候恢复回来
	defer func() {
		if err == nil {
			k.previousContext = strings.TrimSpace(currentContext.String())
		}
	}()

	// 2. 从 clusterConfig 创建 context
	// context 组成：
	// - cluster
	// - redentials
	// - user
	// 三个元素组成 context
	// 2.1 设置 cluster
	// 命令： kubectl config set-cluster <contextName> --server <MasterIP> --insecure-skip-tls-verify=true
	if err := k.Command(
		"kubectl",
		"config",
		"set-cluster",
		contextName,
		"--server",
		config.MasterIP,
		"--insecure-skip-tls-verify=true",
	).Run(); err != nil {
		return err
	}

	// 2.2 设置授权
	if config.Rest.BearerToken != "" {
		// 命令： kubectl config set-credentials <contextName> --token <config.Rest.BearerToken>
		if err := k.Command(
			"kubectl",
			"config",
			"set-credentials",
			contextName,
			"--token",
			config.Rest.BearerToken,
		).Run(); err != nil {
			return err
		}
	} else if config.Rest.CertData != nil && config.Rest.KeyData != nil {
		keyFile := fmt.Sprintf(KeyTempFileFmt, contextName)
		//ginkgo.By("keyFile: " + keyFile)
		certFile := fmt.Sprintf(CertTempFileFmt, contextName)
		//ginkgo.By("certFile: " + certFile)
		if err := os.WriteFile(keyFile, config.Rest.KeyData, os.ModePerm); err != nil {
			return err
		}
		defer func() {
			_ = os.Remove(keyFile)
		}()
		if err := os.WriteFile(certFile, config.Rest.CertData, os.ModePerm); err != nil {
			return err
		}
		defer func() {
			_ = os.Remove(certFile)
		}()

		// 命令：kubectl config set-credentials <contextName> --embed-certs=true --client-key=<keyFile> --client-certificate=<certFile>
		if err := k.Command(
			"kubectl",
			"config",
			"set-credentials",
			contextName,
			"--embed-certs=true",
			"--client-key="+keyFile,
			"--client-certificate="+certFile,
		).Run(); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("counld not find credentials in config or credential method is not supported")
	}

	// 2.3 设置user
	// 命令： kubectl config set-context <contextName> --cluster <contextName> --user=<contextName>
	if err := k.Command(
		"kubectl",
		"config",
		"set-context",
		contextName,
		"--cluster="+contextName,
		"--user="+contextName,
	).Run(); err != nil {
		return err
	}

	// 3. 切换到新创建的 context
	// 命令： kubectl config use-context <contextName>
	if err := k.Command(
		"kubectl",
		"config",
		"use-context",
		contextName,
	).Run(); err != nil {
		return err
	}

	return nil
}

func (k *KubectlConfig) DeleteContext(config *ClusterConfig) error {
	contextName := k.getContextName(config)

	// 1. 删除 cluster
	// 命令: kubectl config delete-cluster <contextName>
	if err := k.Command(
		"kubectl",
		"config",
		"delete-cluster",
		contextName,
	).Run(); err != nil {
		return err
	}

	// 2. 删除 user
	// 命令: kubectl config delete-user <contextName>
	if err := k.Command(
		"kubectl",
		"config",
		"delete-user",
		contextName,
	).Run(); err != nil {
		return err
	}

	// 3. 删除 context
	// 命令： kubectl config delete-context <contextName>
	if err := k.Command(
		"kubectl",
		"config",
		"delete-context",
		contextName,
	).Run(); err != nil {
		return err
	}

	// 4. previousContext 不为空， 恢复到之前的 context
	if k.previousContext != "" {
		// 命令： kubectl config use-context <k.previousContext>
		// 不处理执行错误，因为这里执行错误不影响测试任务
		_ = k.Command(
			"kubectl",
			"config",
			"use-context",
			k.previousContext,
		).Run()
	}

	return nil
}

func (k *KubectlConfig) getContextName(config *ClusterConfig) string {
	return fmt.Sprintf("%s-context", contextNameRegex.ReplaceAllString(strings.ToLower(config.Name), ""))
}

func (k *KubectlConfig) Command(cmd string, args ...string) *exec.Cmd {
	command := exec.Command(cmd, args...)
	command.Stdout = k.config.Stdout
	command.Stderr = k.config.Stderr
	return command
}
