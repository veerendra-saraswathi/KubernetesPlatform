
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// NodeList models the Kubernetes node list JSON
type NodeList struct {
	Items []struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
		Status struct {
			Capacity    map[string]string `json:"capacity"`
			Allocatable map[string]string `json:"allocatable"`
		} `json:"status"`
	} `json:"items"`
}

func main() {
	fmt.Println("🖥️  Kubernetes GPU Detector")
	fmt.Println("==========================")

	// Run kubectl get nodes -o json
	out, err := exec.Command("kubectl", "get", "nodes", "-o", "json").Output()
	if err != nil {
		log.Fatalf("❌ Failed to run kubectl: %v", err)
	}

	// Parse JSON into NodeList
	var nodeList NodeList
	if err := json.Unmarshal(out, &nodeList); err != nil {
		log.Fatalf("❌ JSON parse error: %v", err)
	}

	if len(nodeList.Items) == 0 {
		fmt.Println("⚠️  No nodes found in the cluster")
		return
	}

	// Iterate nodes and check GPU capacity
	for _, node := range nodeList.Items {
		fmt.Printf("\nNode: %s\n", node.Metadata.Name)

		if gpuCount, ok := node.Status.Capacity["nvidia.com/gpu"]; ok {
			fmt.Printf("✅ GPUs detected: %s (Allocatable: %s)\n",
				gpuCount,
				node.Status.Allocatable["nvidia.com/gpu"])
		} else {
			fmt.Println("❌ No GPUs detected on this node")
		}
	}

	fmt.Println("\nGPU check complete.")

}

