package fleet

import (
    "context"
    "fmt"
    "log"

    "k8s.io/client-go/dynamic"
    "k8s.io/client-go/rest"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime/schema"
)

// ListClusters lists all GpuJob resources across all namespaces
func ListClusters() {
    config, err := rest.InClusterConfig()
    if err != nil {
        // fallback to kubeconfig for local testing
        config, err = rest.InClusterConfig()
        if err != nil {
            log.Fatalf("Failed to get Kubernetes config: %v", err)
        }
    }

    dynClient, err := dynamic.NewForConfig(config)
    if err != nil {
        log.Fatalf("Failed to create dynamic client: %v", err)
    }

    gvr := schema.GroupVersionResource{
        Group:    "gpu.yourstartup.com",
        Version:  "v1",
        Resource: "gpujobs",
    }

    list, err := dynClient.Resource(gvr).Namespace("").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        log.Fatalf("Failed to list GpuJobs: %v", err)
    }

    for _, item := range list.Items {
        fmt.Printf("Found GpuJob: %s in namespace %s\n", item.GetName(), item.GetNamespace())
    }
}

