package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// runCommand executes a shell command with a timeout and returns its standard output as a string.
func runCommand(command string, args ...string) (string, error) {
	// Create the command object
	cmd := exec.Command(command, args...)

	// Capture the standard output pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("error creating stdout pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("error starting command: %v", err)
	}

	// Create a channel to wait for the command to finish
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// Set a timeout (e.g., 10 seconds)
	timeout := time.After(10 * time.Second)

	// Read the output from the pipe while command is running
	scanner := bufio.NewScanner(stdout)
	var output strings.Builder
	for scanner.Scan() {
		output.WriteString(scanner.Text() + "\n")
	}

	// Wait for either command completion or timeout
	select {
	case err := <-done:
		if err != nil {
			return output.String(), fmt.Errorf("error waiting for command: %v", err)
		}
	case <-timeout:
		// If timeout occurs, kill the process
		cmd.Process.Kill()
		return "", fmt.Errorf("command timed out after 10 seconds")
	}

	return output.String(), nil
}

// checkNodes checks the status of all nodes in the cluster.
func checkNodes() {
	fmt.Println("=== Node Health Check ===")

	// Run 'kubectl get nodes'
	output, err := runCommand("kubectl", "get", "nodes")
	if err != nil {
		fmt.Printf("❌ Error running kubectl: %v\n", err)
		fmt.Println("Please ensure:")
		fmt.Println("1. A Kubernetes cluster is running")
		fmt.Println("2. kubectl is configured to point to the right cluster")
		fmt.Println("3. Try: kubectl cluster-info to check your connection")
		fmt.Println("")
		fmt.Println("To set up a local cluster:")
		fmt.Println("- Install Minikube: brew install minikube && minikube start")
		fmt.Println("- Or enable Kubernetes in Docker Desktop settings")
		return // Don't fatal error, just return gracefully
	}

	// Parse the output line by line
	scanner := bufio.NewScanner(strings.NewReader(output))
	nodeCount := 0
	readyCount := 0
	notReadyNodes := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		// Skip the header line and empty lines
		if strings.HasPrefix(line, "NAME") || strings.TrimSpace(line) == "" {
			continue
		}
		// Split the line by whitespace. The status is typically in the second column.
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue // Skip malformed lines
		}
		nodeName := fields[0]
		status := fields[1]

		nodeCount++
		// Check if the node is in a 'Ready' state
		if status == "Ready" {
			readyCount++
		} else {
			notReadyNodes = append(notReadyNodes, fmt.Sprintf("%s (%s)", nodeName, status))
		}
	}

	// Print the report
	if nodeCount == 0 {
		fmt.Println("❌ No nodes found in the cluster.")
	} else if readyCount == nodeCount {
		fmt.Printf("✅ All %d node(s) are Ready.\n", nodeCount)
	} else {
		fmt.Printf("⚠️  Cluster Node Status: %d out of %d node(s) are Ready.\n", readyCount, nodeCount)
		fmt.Printf("   Nodes not ready: %s\n", strings.Join(notReadyNodes, ", "))
	}
}

// checkPods checks the status of pods in all namespaces.
func checkPods() {
	fmt.Println("\n=== Pod Health Check (All Namespaces) ===")

	// Run 'kubectl get pods --all-namespaces'
	output, err := runCommand("kubectl", "get", "pods", "-A")
	if err != nil {
		fmt.Printf("❌ Error running kubectl: %v\n", err)
		fmt.Println("Skipping pod check due to kubectl error")
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	totalPods := 0
	pendingPods := 0
	runningPods := 0
	otherPods := 0

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "NAMESPACE") || strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		// The pod status is typically the 4th column in the 'get pods -A' output.
		// Format: NAMESPACE NAME READY STATUS RESTARTS AGE
		if len(fields) < 4 {
			continue
		}
		// namespace := fields[0]
		// podName := fields[1]
		status := fields[3]

		totalPods++

		switch {
		case status == "Running":
			runningPods++
		case status == "Pending":
			pendingPods++
		default:
			otherPods++ // Includes status like "Error", "CrashLoopBackOff", "Completed", etc.
		}
	}

	// Print the report
	fmt.Printf("Total Pods: %d\n", totalPods)
	fmt.Printf("✅ Running: %d\n", runningPods)
	fmt.Printf("⏳ Pending: %d\n", pendingPods)
	fmt.Printf("❓ Other (Error, Completed, etc.): %d\n", otherPods)

	if pendingPods > 0 {
		fmt.Printf("🔍 Note: There are %d pod(s) pending. You may want to investigate.\n", pendingPods)
	}
}

func main() {
	fmt.Println("Cluster Resource Checker Starting...")

	// First, check if kubectl is even available
	_, err := runCommand("kubectl", "version", "--client")
	if err != nil {
		fmt.Printf("❌ kubectl is not available: %v\n", err)
		fmt.Println("Please install kubectl first:")
		fmt.Println("brew install kubectl")
		return
	}
	fmt.Println("✅ kubectl is installed")

	fmt.Println("(Ensure 'kubectl' is configured correctly for your cluster)")
	checkNodes()
	checkPods()
	fmt.Println("\nCheck complete.")
}

