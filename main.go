package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig = flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	var kubecontext = flag.String("context", "default", "context to use")
	var namespace = flag.String("namespace", "default", "namespace of the pod")
	var pod = flag.String("pod", "", "pod name")
	var container = flag.String("container", "", "container name")
	flag.Parse()

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig},
		&clientcmd.ConfigOverrides{CurrentContext: *kubecontext}).ClientConfig()
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	req := clientset.CoreV1().Pods(*namespace).GetLogs(*pod, &v1.PodLogOptions{Container: *container})
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
