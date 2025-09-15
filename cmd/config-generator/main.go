package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AppName     string            `mapstructure:"appName" yaml:"appName"`
	Image       string            `mapstructure:"image" yaml:"image"`
	CPU         string            `mapstructure:"cpu" yaml:"cpu"`
	Memory      string            `mapstructure:"memory" yaml:"memory"`
	GPU         string            `mapstructure:"gpu" yaml:"gpu"`
	GPUType     string            `mapstructure:"gpuType" yaml:"gpuType"`
	Port        int               `mapstructure:"port" yaml:"port"`
	Labels      map[string]string `mapstructure:"labels" yaml:"labels"`
	Environment map[string]string `mapstructure:"environment" yaml:"environment"`
}

func main() {
	// Load config.yaml
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // look in current directory

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	// Print config to console
	fmt.Println("✅ Loaded Configuration:")
	fmt.Printf("App: %s\nImage: %s\nCPU: %s\nMemory: %s\nGPU: %s (%s)\nPort: %d\n",
		cfg.AppName, cfg.Image, cfg.CPU, cfg.Memory, cfg.GPU, cfg.GPUType, cfg.Port)

	// Build Deployment YAML structure
	deployment := map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"name":   cfg.AppName,
			"labels": cfg.Labels,
		},
		"spec": map[string]interface{}{
			"replicas": 1,
			"selector": map[string]interface{}{
				"matchLabels": cfg.Labels,
			},
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": cfg.Labels,
				},
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name":  cfg.AppName,
							"image": cfg.Image,
							"ports": []map[string]interface{}{
								{
									"containerPort": cfg.Port,
								},
							},
							"resources": map[string]interface{}{
								"requests": map[string]string{
									"cpu":    cfg.CPU,
									"memory": cfg.Memory,
								},
								"limits": map[string]string{
									"cpu":    cfg.CPU,
									"memory": cfg.Memory,
								},
							},
							"env": func() []map[string]string {
								var envs []map[string]string
								for k, v := range cfg.Environment {
									envs = append(envs, map[string]string{
										"name":  k,
										"value": v,
									})
								}
								return envs
							}(),
						},
					},
				},
			},
		},
	}

	// Marshal to YAML
	yamlData, err := yaml.Marshal(deployment)
	if err != nil {
		log.Fatalf("Error marshaling Deployment YAML: %v", err)
	}

	// Print generated Deployment YAML
	fmt.Println("\n📦 Generated Deployment YAML:")
	fmt.Println(string(yamlData))

	// Save to file
	if err := os.WriteFile("deployment.yaml", yamlData, 0644); err != nil {
		log.Fatalf("Error writing deployment.yaml: %v", err)
	}
	fmt.Println("\n💾 Deployment YAML written to deployment.yaml")
}

