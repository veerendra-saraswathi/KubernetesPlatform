/*
Phase 2+: Pod Creator CLI Enhancements

Next steps / TODOs:
- Extend the CLI to delete Pods too.
- Add more resource handling:
    * CPU/GPU/Memory requests & limits
    * Multi-container Pods
    * Configurable environment variables
- Dry-run / simulation mode:
    * Print the YAML instead of applying it, useful for AI/ML workloads and clusters without GPUs.
- Integration:
    * Combine health-server + pod-creator + GPU lister to have a mini self-monitoring platform.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Command-line flags
	podName := flag.String("name", "my-pod", "Name of the Pod")
	image := flag.String("image", "nginx", "Container image")
	namespace := flag.String("namespace", "default", "Namespace for the Pod")
	flag.Parse()

	// Build kubeconfig path
	kubeconfig := filepath.Join(homeDir(), ".kube", "config")

	// Create Kubernetes clientset
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Define Pod object
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *podName,
			Namespace: *namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  *podName,
					Image: *image,
				},
			},
		},
	}

	// Create the Pod
	createdPod, err := clientset.CoreV1().Pods(*namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("✅ Pod '%s' created in namespace '%s'\n", createdPod.Name, *namespace)

	// Wait for Pod to be Running and Ready
	fmt.Println("⏳ Waiting for Pod to be Running and Ready...")
	for {
		p, err := clientset.CoreV1().Pods(*namespace).Get(context.TODO(), *podName, metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}

		ready := false
		for _, cond := range p.Status.Conditions {
			if cond.Type == v1.PodReady && cond.Status == v1.ConditionTrue {
				ready = true
				break
			}
		}

		if ready {
			fmt.Printf("🟢 Pod '%s' and all containers are Ready\n", *podName)
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Fetch and print Pod logs (last 20 lines)
	fmt.Println("📜 Fetching logs (last 20 lines)...")
	logOptions := &v1.PodLogOptions{
		TailLines: int64Ptr(20),
	}
	req := clientset.CoreV1().Pods(*namespace).GetLogs(*podName, logOptions)
	logs, err := req.Stream(context.TODO())
	if err != nil {
		panic(err.Error())
	}
	defer logs.Close()
	buf := make([]byte, 2000)
	n, _ := logs.Read(buf)
	fmt.Println(string(buf[:n]))
}

// Helper function to get home directory
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return "."
}

// Helper to get pointer to int64
func int64Ptr(i int64) *int64 {
	return &i
}

