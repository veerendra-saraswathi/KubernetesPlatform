/*
Phase 2: Native GPU Node Lister using client-go

Next steps / TODOs:
- Test with a real GPU node later — it will automatically pick up nvidia.com/gpu, amd.com/gpu, etc.
- Extend this CLI:
    * Print Allocatable vs Capacity for CPU/GPU/Memory.
    * Output JSON or YAML for automation.
    * Combine with your health-server to check node & service health.
*/

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Build kubeconfig path from HOME directory
	kubeconfig := filepath.Join(homeDir(), ".kube", "config")

	// Create Kubernetes clientset from kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// List all nodes in the cluster
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Found %d node(s) in cluster:\n\n", len(nodes.Items))

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

		fmt.Println("---------------------------")
	}

	fmt.Printf("Total GPU Nodes: %d | Total GPUs: %d\n", totalGPUNodes, totalGPUs)
}

// helper function to get the home directory
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return "."
}

