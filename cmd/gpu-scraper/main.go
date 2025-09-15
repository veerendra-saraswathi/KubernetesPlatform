/*
Phase 3: GPU Usage Scraper & Notifier (Pre-Operator)

Goals / TODOs:
- Periodically check allocatable vs allocated GPUs on each node.
- Print a warning if GPU usage > 90%.
- Learn simple controller logic before writing full Kubernetes operators.

Extended Functionality / Next Steps:
- Send Slack/email alerts when GPU usage exceeds threshold.
- Track multiple GPU types per node.
- Output JSON/YAML for integration with other monitoring tools.
- Extend to metrics server integration or AI/ML monitoring dashboards.
*/

package main

import (
    "context"
    "encoding/json"
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

    // GPU resources to track
    gpuResources := []string{"nvidia.com/gpu", "amd.com/gpu", "intel.com/gpu"}

    fmt.Println("🚀 Starting GPU Usage Scraper...")
    ticker := time.NewTicker(10 * time.Second) // adjust interval as needed
    defer ticker.Stop()

    for range ticker.C {
        nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
        if err != nil {
            fmt.Println("❌ Error listing nodes:", err)
            continue
        }

        results := make(map[string]map[string]interface{}) // for JSON/YAML output

        for _, node := range nodes.Items {
            fmt.Printf("\nNode: %s\n", node.Name)
            nodeResults := make(map[string]interface{})

            for _, gpuRes := range gpuResources {
                allocatable, ok := node.Status.Allocatable[v1.ResourceName(gpuRes)]
                if !ok || allocatable.Value() == 0 {
                    continue
                }

                // Calculate total GPUs currently allocated to Pods
                pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
                if err != nil {
                    fmt.Println("❌ Error listing pods:", err)
                    continue
                }

                usedGPUs := int64(0)
                for _, pod := range pods.Items {
                    if pod.Spec.NodeName != node.Name {
                        continue
                    }
                    for _, container := range pod.Spec.Containers {
                        if val, ok := container.Resources.Requests[v1.ResourceName(gpuRes)]; ok {
                            usedGPUs += val.Value()
                        }
                    }
                }

                allocVal := allocatable.Value()
                usagePercent := float64(usedGPUs) / float64(allocVal) * 100
                fmt.Printf("%s: Allocatable=%d, Used=%d, Usage=%.2f%%\n", gpuRes, allocVal, usedGPUs, usagePercent)

                // Populate results for JSON/YAML output
                nodeResults[string(gpuRes)] = map[string]interface{}{
                    "allocatable": allocVal,
                    "used":        usedGPUs,
                    "usagePercent": usagePercent,
                }

                // Warning & alert
                if usagePercent >= 90 {
                    fmt.Printf("⚠️  GPU usage above 90%% on node %s for %s\n", node.Name, gpuRes)
                    // TODO: Send Slack/email alert here
                }
            }

            results[node.Name] = nodeResults
        }

        // Optional: Output JSON for automation
        outputJSON, _ := json.MarshalIndent(results, "", "  ")
        fmt.Println("\n📊 GPU Usage JSON Output:")
        fmt.Println(string(outputJSON))
        // TODO: Write to file or integrate with monitoring system
    }
}

// helper function to get home directory
func homeDir() string {
    if h := os.Getenv("HOME"); h != "" {
        return h
    }
    return "."
}

