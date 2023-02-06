package kube

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

var (
	configPath string
)

type kubeClient struct {
	client *kubernetes.Clientset
}

func NewKubeClient() (*kubeClient, error) {
	fmt.Println("new kube client")
	var config *rest.Config
	var err error
	if configPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			return nil, err
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &kubeClient{
		client: clientset,
	}, nil
}

func (k *kubeClient) Update(namespace, deploymentName string) error {
	fmt.Println("new kube client")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	deployment, err := k.client.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	containers := &deployment.Spec.Template.Spec.Containers
	fmt.Println(containers)
	for i := range *containers {
		c := *containers
		for j := range c[i].Env {
			if c[i].Env[j].Name == "timestep" {
				c[i].Env[j].Value = time.Now().String()
			}
		}
	}
	_, err = k.client.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func init() {
	flag.StringVar(&configPath, "kubeconfig", "", "kube config, eg: -kubeconfig xxx")
}
