package dashboard

import (
    "context"
    "fmt"
    "k8s.io/client-go/dynamic"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/apimachinery/pkg/runtime/schema"
)

func ShowClusterHealth(kubeconfigPath string) error {
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
    if err != nil {
        return fmt.Errorf("failed to load kubeconfig: %v", err)
    }

    dynClient, err := dynamic.NewForConfig(config)
    if err != nil {
        return fmt.Errorf("failed to create dynamic client: %v", err)
    }

    clustersRes := schema.GroupVersionResource{
        Group:    "cluster.yourstartup.com",
        Version:  "v1",
        Resource: "managedclusters",
    }

    list, err := dynClient.Resource(clustersRes).Namespace("").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        return fmt.Errorf("failed to list clusters: %v", err)
    }

    fmt.Println("Clusters found:", len(list.Items))
    for _, cluster := range list.Items {
        fmt.Println(" -", cluster.GetName())
    }

    return nil
}

