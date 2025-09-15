/*
Phase 3: Job Status Poller

Goals / TODOs:
- Poll a Kubernetes Job periodically after submission.
- Print Active, Succeeded, or Failed status.
- Can be combined with Job Submitter for automation.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	jobName := flag.String("job", "", "Name of the Job to poll")
	namespace := flag.String("namespace", "default", "Namespace of the Job")
	flag.Parse()

	if *jobName == "" {
		fmt.Println("❌ Please provide a Job name using -job")
		os.Exit(1)
	}

	kubeconfig := filepath.Join(homeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("⏳ Polling Job '%s' in namespace '%s'...\n", *jobName, *namespace)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		job, err := clientset.BatchV1().Jobs(*namespace).Get(context.TODO(), *jobName, metav1.GetOptions{})
		if err != nil {
			fmt.Println("❌ Error fetching Job:", err)
			continue
		}

		active := job.Status.Active
		succeeded := job.Status.Succeeded
		failed := job.Status.Failed

		fmt.Printf("Job Status - Active: %d, Succeeded: %d, Failed: %d\n", active, succeeded, failed)

		if succeeded > 0 {
			fmt.Println("✅ Job completed successfully")
			break
		} else if failed > 0 {
			fmt.Println("❌ Job failed")
			break
		}
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return "."
}

