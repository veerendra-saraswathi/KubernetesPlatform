/*
Phase 4: CRD Installer

Goals / TODOs:
- Read a YAML file containing a CustomResourceDefinition (CRD).
- Use apiextensions-apiserver client to create the CRD in the cluster.
- Learn the first step towards building a Kubernetes Operator.
*/

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	apixv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apixclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/util/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var crdFile string
	flag.StringVar(&crdFile, "file", "", "Path to CRD YAML file")
	flag.Parse()

	if crdFile == "" {
		fmt.Println("❌ Please provide a CRD YAML file using -file flag")
		os.Exit(1)
	}

	// Read CRD YAML
	data, err := ioutil.ReadFile(crdFile)
	if err != nil {
		panic(fmt.Errorf("failed to read CRD file: %v", err))
	}

	// Build kubeconfig path from HOME if not in KUBECONFIG
	kubeconfig := filepath.Join(homeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(fmt.Errorf("failed to build kubeconfig: %v", err))
	}

	// Create API extensions client
	clientset, err := apixclient.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("failed to create API extensions client: %v", err))
	}

	// Decode YAML into a CRD object
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), 100)
	var crd apixv1.CustomResourceDefinition
	if err := decoder.Decode(&crd); err != nil {
		panic(fmt.Errorf("failed to decode CRD YAML: %v", err))
	}

	// Check if CRD already exists
	_, err = clientset.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), crd.Name, metav1.GetOptions{})
	if err == nil {
		fmt.Printf("ℹ️  CRD '%s' already exists, skipping creation\n", crd.Name)
		return
	}

	// Create CRD
	created, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(), &crd, metav1.CreateOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to create CRD: %v", err))
	}

	fmt.Printf("✅ CRD '%s' created successfully\n", created.Name)
}

// helper function to get home directory
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return "."
}

