package main

import (
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/cache"
)

func main() {
	kubeconfig := filepath.Join(homeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create a shared informer factory with 0 resync period
	factory := informers.NewSharedInformerFactory(clientset, 0)

	// Get a namespace informer
	namespaceInformer := factory.Core().V1().Namespaces().Informer()

	// Add event handlers
	namespaceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ns := obj.(*v1.Namespace)
			fmt.Printf("🟢 Namespace created: %s\n", ns.Name)
		},
		DeleteFunc: func(obj interface{}) {
			ns := obj.(*v1.Namespace)
			fmt.Printf("🔴 Namespace deleted: %s\n", ns.Name)
		},
	})

	stopCh := make(chan struct{})
	defer close(stopCh)

	fmt.Println("Starting Namespace Watcher...")
	go namespaceInformer.Run(stopCh)

	// Wait forever
	select {}
}

// helper function to get home directory
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return "."
}

