package fleet

import (
    "context"
    "fmt"

    "k8s.io/client-go/dynamic"
    "k8s.io/client-go/tools/clientcmd"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListClusters() {
    // load kubeconfig
    config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
    if err != nil {
        panic(err)
    }

    dynClient, err := dynamic.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    // replace group/resource with your ManagedCluster CRD
    gvr := dynamic.Resource{
        Group:    "cluster.yourstartup.com",
        Version:  "v1",
        Resource: "managedclusters",
    }

    clusters, err := dynClient.Resource(gvr).Namespace("").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        panic(err)
    }

    for _, c := range clusters.Items {
        fmt.Printf("Found cluster: %s\n", c.GetName())
    }
}

