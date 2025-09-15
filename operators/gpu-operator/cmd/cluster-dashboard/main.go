package main

import (
    "fmt"
    "os"
    "yourstartup.com/gpu-operator/internal/controller/dashboard"
)

func main() {
    kubeconfig := os.Getenv("KUBECONFIG")
    if kubeconfig == "" {
        kubeconfig = os.ExpandEnv("$HOME/.kube/config")
    }

    err := dashboard.ShowClusterHealth(kubeconfig)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
}

