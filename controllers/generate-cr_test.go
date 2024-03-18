package controllers

import (
	"fmt"
	myAppv1 "github.com/huangxinzhong/huangxinzhong-deployment/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"reflect"
	"testing"
)

func readFile(fileName string) []byte {
	content, err := os.ReadFile(fmt.Sprintf("testdata/%s", fileName))
	if err != nil {
		panic(err)
	}
	return content
}

func newHxzDeployment(fileName string) *myAppv1.HxzDeployment {
	content := readFile(fileName)
	md := new(myAppv1.HxzDeployment)
	if err := yaml.Unmarshal(content, md); err != nil {
		panic(err)
	}
	return md
}

func newDeployment(fileName string) *appsv1.Deployment {
	content := readFile(fileName)
	md := new(appsv1.Deployment)
	if err := yaml.Unmarshal(content, md); err != nil {
		panic(err)
	}
	return md
}

func newService(fileName string) *corev1.Service {
	content := readFile(fileName)
	md := new(corev1.Service)
	if err := yaml.Unmarshal(content, md); err != nil {
		panic(err)
	}
	return md
}

func newServiceNP(fileName string) *corev1.Service {
	content := readFile(fileName)
	md := new(corev1.Service)
	if err := yaml.Unmarshal(content, md); err != nil {
		panic(err)
	}
	return md
}

func newIngress(fileName string) *networkv1.Ingress {
	content := readFile(fileName)
	md := new(networkv1.Ingress)
	if err := yaml.Unmarshal(content, md); err != nil {
		panic(err)
	}
	return md
}

func TestNewDeployment(t *testing.T) {
	type args struct {
		md *myAppv1.HxzDeployment
	}
	tests := []struct {
		name    string             // 测试用例名称
		args    args               // 测试函数的参数
		want    *appsv1.Deployment // 期望结果
		wantErr bool               // 进行测试时函数是否需要出错
	}{
		{
			name: "测试使用 ingress mode 时候，生成 Deployment 资源",
			args: args{
				md: newHxzDeployment("hxz-ingress-cr.yaml"),
			},
			want:    newDeployment("hxz-ingress-deployment-expect.yaml"),
			wantErr: false,
		},
		{
			name: "测试使用 nodeport mode 时候，生成 Deployment 资源",
			args: args{
				md: newHxzDeployment("hxz-nodeport-cr.yaml"),
			},
			want:    newDeployment("hxz-nodeport-deployment-expect.yaml"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDeployment(tt.args.md)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDeployment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeployment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIngress(t *testing.T) {
	type args struct {
		md *myAppv1.HxzDeployment
	}
	tests := []struct {
		name    string
		args    args
		want    *networkv1.Ingress
		wantErr bool
	}{
		{
			name: "测试使用 ingress mode 时候， 生成 ingress 资源",
			args: args{
				md: newHxzDeployment("hxz-ingress-cr.yaml"),
			},
			want:    newIngress("hxz-ingress-ingress-expect.yaml"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewIngress(tt.args.md)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewIngress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIngress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewService(t *testing.T) {
	type args struct {
		md *myAppv1.HxzDeployment
	}
	tests := []struct {
		name    string
		args    args
		want    *corev1.Service
		wantErr bool
	}{
		{
			name: "测试使用 ingress mode时候生成 service 资源",
			args: args{
				md: newHxzDeployment("hxz-ingress-cr.yaml"),
			},
			want:    newService("hxz-ingress-service-expect.yaml"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewService(tt.args.md)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewServiceNodePort(t *testing.T) {
	type args struct {
		md *myAppv1.HxzDeployment
	}
	tests := []struct {
		name    string
		args    args
		want    *corev1.Service
		wantErr bool
	}{
		{
			name: "测试使用 nodeport mode 时候生成 nodeport 类型的 service 资源",
			args: args{
				md: newHxzDeployment("hxz-nodeport-cr.yaml"),
			},
			want:    newService("hxz-nodeport-service-expect.yaml"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewServiceNodePort(tt.args.md)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServiceNodePort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServiceNodePort() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTemplate(t *testing.T) {
	type args struct {
		md           *myAppv1.HxzDeployment
		templateName string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTemplate(tt.args.md, tt.args.templateName)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTemplate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
