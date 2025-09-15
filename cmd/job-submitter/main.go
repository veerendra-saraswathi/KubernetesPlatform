/*
Phase 3: Simple Job Submitter

Goals / TODOs:
- Takes a Kubernetes Job YAML file and submits it to the cluster.
- Automates benchmark or CUDA test jobs.
- Can be extended to handle multiple jobs, GPU resource checks, or dry-run mode.
*/

package main

import (
    "context"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"

    batchv1 "k8s.io/api/batch/v1"
    "k8s.io/apimachinery/pkg/util/yaml"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    // Command-line flags
    yamlFile := flag.String("file", "", "Path to Kubernetes Job YAML file")
    namespace := flag.String("namespace", "default", "Namespace to submit the Job")
    flag.Parse()

    if *yamlFile == "" {
        fmt.Println("❌ Please provide the path to a Job YAML file using -file")
        return
    }

    // Read YAML file
    data, err := ioutil.ReadFile(*yamlFile)
    if err != nil {
        panic(err)
    }

    // Unmarshal YAML into Job object
    job := &batchv1.Job{}
    if err := yaml.Unmarshal(data, job); err != nil {
        panic(err)
    }

    // Build kubeconfig path
    kubeconfig := filepath.Join(homeDir(), ".kube", "config")

    // Create Kubernetes clientset
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        panic(err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    // Submit Job
    createdJob, err := clientset.BatchV1().Jobs(*namespace).Create(context.TODO(), job, metav1.CreateOptions{})
    if err != nil {
        panic(err)
    }

    fmt.Printf("✅ Job '%s' submitted in namespace '%s'\n", createdJob.Name, *namespace)
}

// helper function to get home directory
func homeDir() string {
    if h := os.Getenv("HOME"); h != "" {
        return h
    }
    return "."
}

