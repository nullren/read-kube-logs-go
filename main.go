package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig = flag.String("c", "/etc/kube/config", "location of kubeconfig to read")
var namespace = flag.String("n", "default", "namespace of the pod")
var pod = flag.String("p", "", "pod name")

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	req := clientset.CoreV1().Pods(*namespace).GetLogs(*pod, &v1.PodLogOptions{})
	podLogs, err := req.Stream(ctx)
	if err != nil {
		panic(err)
	}
	defer podLogs.Close()

	scanner := bufio.NewScanner(podLogs)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
}
