/*
Phase 2: Combined Cluster Monitor CLI
- GPU Node Lister using client-go
- Namespace Watcher using SharedInformer
Next steps / TODOs:
- Print Allocatable vs Capacity for CPU/GPU/Memory
- Output JSON or YAML for automation
- Combine with health-server to check node & service health
- Extend watchers to Pods, Deployments, and Services
*/

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// ------------------- Kubeconfig -------------------
	kubeconfig := filepath.Join(homeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Starting Cluster Monitor...")

	// ---------------- GPU Node Lister ----------------
	go func() {
		for {
			fmt.Println("\n=== GPU Node Check ===")
			checkGPUNodes(clientset)
			time.Sleep(30 * time.Second) // TODO: make interval configurable
		}
	}()

	// ---------------- Namespace Watcher ----------------
	factory := informers.NewSharedInformerFactory(clientset, 0)
	namespaceInformer := factory.Core().V1().Namespaces().Informer()

	namespaceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ns := obj.(*v1.Namespace)
			fmt.Printf("🟢 Namespace created: %s\n", ns.Name)
			// TODO: trigger automated workflow for new namespace
		},
		DeleteFunc: func(obj interface{}) {
			ns := obj.(*v1.Namespace)
			fmt.Printf("🔴 Namespace deleted: %s\n", ns.Name)
			// TODO: clean up associated resources if needed
		},
	})

	stopCh := make(chan struct{})
	defer close(stopCh)
	go namespaceInformer.Run(stopCh)

	// ---------------- Keep main alive ----------------
	select {} // TODO: add graceful shutdown handling
}

// ---------------- GPU Node Lister Function ----------------
func checkGPUNodes(clientset *kubernetes.Clientset) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing nodes:", err)
		return
	}

	gpuResources := []string{"nvidia.com/gpu", "amd.com/gpu", "intel.com/gpu"}
	totalGPUNodes := 0
	totalGPUs := 0

	for _, node := range nodes.Items {
		fmt.Printf("Node: %s\n", node.Name)

		// Check if node is Ready
		ready := false
		for _, cond := range node.Status.Conditions {
			if cond.Type == v1.NodeReady && cond.Status == v1.ConditionTrue {
				ready = true
				break
			}
		}
		status := "❌ Not Ready"
		if ready {
			status = "✅ Ready"
		}
		fmt.Println("Status:", status)

		// Check for GPU resources
		foundGPU := false
		for _, res := range gpuResources {
			if count, ok := node.Status.Capacity[v1.ResourceName(res)]; ok {
				if count.Value() > 0 {
					fmt.Printf("GPU Detected: %s x %d\n", res, count.Value())
					foundGPU = true
					totalGPUNodes++
					totalGPUs += int(count.Value())
				}
			}
		}

		if !foundGPU {
			fmt.Println("No GPUs detected on this node.")
		}

		// TODO: Print Allocatable vs Capacity for CPU/Memory/GPU
		fmt.Println("---------------------------")
	}
	fmt.Printf("Total GPU Nodes: %d | Total GPUs: %d\n", totalGPUNodes, totalGPUs)
}

// ---------------- Helper ----------------
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return "."
}

