package controllers

import (
	"bytes"
	"fmt"
	"text/template"

	"k8s.io/apimachinery/pkg/util/yaml"

	myAppv1 "github.com/huangxinzhong/huangxinzhong-deployment/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
)

func parseTemplate(md *myAppv1.HxzDeployment, templateName string) ([]byte, error) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("controllers/templates/%s.yaml", templateName))
	if err != nil {
		return nil, err
	}
	b := &bytes.Buffer{}
	if err := tmpl.Execute(b, md); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func NewDeployment(md *myAppv1.HxzDeployment) (*appsv1.Deployment, error) {
	content, err := parseTemplate(md, "deployment")
	if err != nil {
		return nil, err
	}

	deploy := new(appsv1.Deployment)
	if err := yaml.Unmarshal(content, deploy); err != nil {
		return nil, err
	}

	return deploy, nil
}

func NewIngress(md *myAppv1.HxzDeployment) (*networkv1.Ingress, error) {
	content, err := parseTemplate(md, "ingress")
	if err != nil {
		return nil, err
	}

	ig := new(networkv1.Ingress)
	if err := yaml.Unmarshal(content, ig); err != nil {
		return nil, err
	}

	return ig, nil
}

func NewService(md *myAppv1.HxzDeployment) (*corev1.Service, error) {
	content, err := parseTemplate(md, "service")
	if err != nil {
		return nil, err
	}

	svc := new(corev1.Service)
	if err := yaml.Unmarshal(content, svc); err != nil {
		return nil, err
	}

	return svc, nil
}

func NewServiceNodePort(md *myAppv1.HxzDeployment) (*corev1.Service, error) {
	content, err := parseTemplate(md, "service-nodeport")
	if err != nil {
		return nil, err
	}

	svc := new(corev1.Service)
	if err := yaml.Unmarshal(content, svc); err != nil {
		return nil, err
	}

	return svc, nil
}
