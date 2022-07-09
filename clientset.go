package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := "D:/code/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, v := range pods.Items {
		fmt.Printf("%v %v\n", v.ObjectMeta.Name, v.ObjectMeta.Namespace)
	}
	deps, err := clientset.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, v := range deps.Items {
		fmt.Printf("%v %v\n", v.ObjectMeta.Name, v.ObjectMeta.Namespace)
	}
}
