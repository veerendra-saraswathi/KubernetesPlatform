/*
Phase 3: Node Tainter for GPU Nodes

Goals / TODOs:
- Find all nodes with GPUs.
- Add taint gpu=true:NoSchedule if not already present.
- Ensures only pods tolerating this taint get scheduled on GPU nodes.
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
    kubeconfig := filepath.Join(homeDir(), ".kube", "config")

    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        panic(err.Error())
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }

    gpuResources := []string{"nvidia.com/gpu", "amd.com/gpu", "intel.com/gpu"}

    nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        panic(err.Error())
    }

    for _, node := range nodes.Items {
        hasGPU := false
        for _, res := range gpuResources {
            if val, ok := node.Status.Capacity[v1.ResourceName(res)]; ok && val.Value() > 0 {
                hasGPU = true
                break
            }
        }

        if hasGPU {
            fmt.Println("GPU Node found:", node.Name)
            // Check if taint already exists
            taintExists := false
            for _, t := range node.Spec.Taints {
                if t.Key == "gpu" && t.Value == "true" && t.Effect == v1.TaintEffectNoSchedule {
                    taintExists = true
                    break
                }
            }

            if !taintExists {
                fmt.Println("Adding taint gpu=true:NoSchedule to node", node.Name)
                taint := v1.Taint{
                    Key:    "gpu",
                    Value:  "true",
                    Effect: v1.TaintEffectNoSchedule,
                }
                node.Spec.Taints = append(node.Spec.Taints, taint)
                _, err := clientset.CoreV1().Nodes().Update(context.TODO(), &node, metav1.UpdateOptions{})
                if err != nil {
                    fmt.Println("Error adding taint:", err)
                } else {
                    fmt.Println("Taint added successfully")
                }
            } else {
                fmt.Println("Node already tainted:", node.Name)
            }
        }
    }
}

func homeDir() string {
    if h := os.Getenv("HOME"); h != "" {
        return h
    }
    return "."
}

