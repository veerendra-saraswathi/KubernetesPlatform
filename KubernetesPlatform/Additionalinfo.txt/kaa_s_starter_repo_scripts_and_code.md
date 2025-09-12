# KaaS Starter Repo — Multi-file Code Dump

This document contains a starter set of scripts and source files to kickstart your **Kubernetes-as-a-Service (KaaS)** backend-first project. It includes Bash scripts for environment bootstrap, a Go backend skeleton that talks to Cluster API (CAPI), a basic Dockerfile/Makefile, and a small Ruby CLI helper. Use this as a reference scaffold — adapt, harden, and extend for production.

---

## Repository layout

```
kaas-starter/
├── README.md
├── Makefile
├── bootstrap/
│   ├── 00-check-prereqs.sh
│   ├── 01-install-clusterctl.sh
│   ├── 02-init-capi-docker.sh
├── scripts/
│   ├── build-and-push.sh
│   └── local-kind-setup.sh
├── backend/
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── handlers.go
│   ├── capi_client.go
│   └── Dockerfile
├── cli/
│   └── clusterctl_helper.rb
└── docs/
    └── quickstart.md
```

---

> NOTE: These files are scaffolding and educational. They assume you will run in a secure dev environment. Secrets, credentials or production hardening (RBAC, TLS, audit logging) must be added before real use.

---

## File: `bootstrap/00-check-prereqs.sh`

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "==> Checking prerequisites for KaaS dev environment"

command -v docker >/dev/null 2>&1 || { echo "docker is required. Install Docker or use rootless docker"; exit 1; }
command -v kubectl >/dev/null 2>&1 || { echo "kubectl is required. Install kubectl"; exit 1; }
command -v clusterctl >/dev/null 2>&1 || { echo "clusterctl is required. Run bootstrap/01-install-clusterctl.sh"; exit 1; }

if [[ $(uname -s) == "Darwin" ]]; then
  echo "Detected macOS host"
fi

echo "Prereqs present. Docker, kubectl and clusterctl found."
```

---

## File: `bootstrap/01-install-clusterctl.sh`

```bash
#!/usr/bin/env bash
set -euo pipefail

# Installs clusterctl on Linux or macOS

TARGET_BIN=${TARGET_BIN:-/usr/local/bin}
CAPI_VERSION=${CAPI_VERSION:-v1.6.0}

echo "Installing clusterctl ${CAPI_VERSION} to ${TARGET_BIN}"

curl -L "https://github.com/kubernetes-sigs/cluster-api/releases/download/${CAPI_VERSION}/clusterctl-${CAPI_VERSION}-$(uname | tr '[:upper:]' '[:lower:]')-amd64" -o /tmp/clusterctl
chmod +x /tmp/clusterctl
sudo mv /tmp/clusterctl ${TARGET_BIN}/clusterctl

clusterctl version
```

---

## File: `bootstrap/02-init-capi-docker.sh`

```bash
#!/usr/bin/env bash
set -euo pipefail

# Initialize CAPI with Docker provider in a local dev environment.
# This uses clusterctl to bootstrap the management cluster.

CAPI_VERSION=${CAPI_VERSION:-v1.6.0}
NAMESPACE=${NAMESPACE:-capi-system}

echo "Initializing Cluster API components (docker provider)"

# Ensure Docker provider images available; this script assumes docker daemon is running
clusterctl init --infrastructure docker

echo "Cluster API initialized with docker provider."

# show pods
kubectl get pods -n capi-system --show-labels || true
```

---

## File: `scripts/local-kind-setup.sh`

```bash
#!/usr/bin/env bash
set -euo pipefail

# Creates a local kind cluster to be managed by CAPI for demo & testing

CLUSTER_NAME=${CLUSTER_NAME:-capi-mgmt}

cat <<EOF | kind create cluster --name ${CLUSTER_NAME} --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
EOF

kubectl cluster-info --context kind-${CLUSTER_NAME}

echo "Kind cluster ${CLUSTER_NAME} created. Use clusterctl init --infrastructure docker to initialize CAPI."
```

---

## File: `scripts/build-and-push.sh`

```bash
#!/usr/bin/env bash
set -euo pipefail

# Build Docker image for backend and push to registry
IMAGE=${IMAGE:-localhost:5000/kaas-backend:dev}

# optionally build with docker buildx for multi-arch
docker build -t ${IMAGE} ./backend

echo "Pushing ${IMAGE}"
docker push ${IMAGE}

echo "Done"
```

---

## File: `backend/go.mod`

```go
module github.com/yourorg/kaas-backend

go 1.21

require (
    k8s.io/apimachinery v0.27.10
    sigs.k8s.io/controller-runtime v0.14.0
    sigs.k8s.io/cluster-api/api/v1beta1 v0.0.0-20240201-000000000000
)
```

> Note: The cluster-api go module version here is illustrative; check latest versions and replace.

---

## File: `backend/capi_client.go`

```go
package main

import (
    "context"
    "fmt"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/client/config"
    capi "sigs.k8s.io/cluster-api/api/v1beta1"
)

// Global client; simple example for POC. In production inject via constructors and use proper contexts
var k8sClient client.Client

func initK8sClient() error {
    cfg, err := config.GetConfig()
    if err != nil {
        return fmt.Errorf("unable to get kubeconfig: %w", err)
    }
    k8sClient, err = client.New(cfg, client.Options{})
    if err != nil {
        return fmt.Errorf("unable to create k8s client: %w", err)
    }
    return nil
}

// listClusters returns Cluster CRs in management cluster namespace
func listClusters(ctx context.Context, namespace string) (*capi.ClusterList, error) {
    var list capi.ClusterList
    if err := k8sClient.List(ctx, &list, &client.ListOptions{Namespace: namespace}); err != nil {
        return nil, err
    }
    return &list, nil
}

// createSimpleCluster is a tiny example: creates a Cluster CR that references a DockerCluster infrastructure
func createSimpleCluster(ctx context.Context, name, namespace string) (*capi.Cluster, error) {
    cluster := &capi.Cluster{
        ObjectMeta: metav1.ObjectMeta{
            Name:      name,
            Namespace: namespace,
        },
        Spec: capi.ClusterSpec{
            InfrastructureRef: &capi.ObjectReference{
                Kind:       "DockerCluster",
                APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
                Name:       name,
            },
        },
    }
    if err := k8sClient.Create(ctx, cluster); err != nil {
        return nil, err
    }
    return cluster, nil
}
```

---

## File: `backend/handlers.go`

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
)

// Minimal HTTP handlers that call into capi_client.go

func registerHandlers(r *mux.Router) {
    r.HandleFunc("/healthz", healthHandler).Methods("GET")
    r.HandleFunc("/v1/clusters", createClusterHandler).Methods("POST")
    r.HandleFunc("/v1/clusters", listClustersHandler).Methods("GET")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "ok")
}

type createClusterRequest struct {
    Name      string `json:"name"`
    Namespace string `json:"namespace"`
}

func createClusterHandler(w http.ResponseWriter, r *http.Request) {
    var req createClusterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    ctx := context.Background()
    if req.Namespace == "" {
        req.Namespace = "default"
    }
    cluster, err := createSimpleCluster(ctx, req.Name, req.Namespace)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cluster)
}

func listClustersHandler(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    namespace := r.URL.Query().Get("namespace")
    if namespace == "" { namespace = "default" }
    list, err := listClusters(ctx, namespace)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(list)
}
```

---

## File: `backend/main.go`

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
)

func main() {
    // Init k8s client
    if err := initK8sClient(); err != nil {
        log.Fatalf("failed to init k8s client: %v", err)
    }

    r := mux.NewRouter()
    registerHandlers(r)

    port := os.Getenv("PORT")
    if port == "" { port = "8080" }

    srv := &http.Server{
        Addr:    ":" + port,
        Handler: r,
    }

    go func() {
        <-context.Background().Done()
        _ = srv.Shutdown(context.Background())
    }()

    fmt.Printf("Starting API server on %s\n", srv.Addr)
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("server error: %v", err)
    }
}
```

---

## File: `backend/Dockerfile`

```dockerfile
FROM golang:1.21-bullseye AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/kaas-backend ./

FROM gcr.io/distroless/static
COPY --from=builder /out/kaas-backend /kaas-backend
ENTRYPOINT ["/kaas-backend"]
```

---

## File: `Makefile`

```makefile
.PHONY: build docker-run test fmt vet

build:
	cd backend && go build -o kaas-backend .

docker-build:
	docker build -t localhost:5000/kaas-backend:dev ./backend

docker-run:
	docker run --rm -p 8080:8080 \
	  -v ${HOME}/.kube/config:/root/.kube/config:ro \
	  localhost:5000/kaas-backend:dev

fmt:
	gofmt -w ./backend

test:
	# placeholder for unit tests
	echo "no tests yet"
```

---

## File: `cli/clusterctl_helper.rb`

```ruby
#!/usr/bin/env ruby
# Minimal Ruby helper that shells out to clusterctl and kubectl for convenience
# Use for quick demos; in production build proper Ruby/Gems or Go CLI

require 'json'
require 'open3'

def run(cmd)
  puts "> #{cmd}"
  stdout, stderr, status = Open3.capture3(cmd)
  unless status.success?
    puts "ERROR: #{stderr}"
    exit 1
  end
  stdout
end

def list_clusters
  puts run('kubectl get clusters -A -o json')
end

def create_demo_cluster(name)
  # example quick-and-dirty yaml - in practice use clusterctl config
  yaml = <<~YAML
  apiVersion: cluster.x-k8s.io/v1beta1
  kind: Cluster
  metadata:
    name: #{name}
    namespace: default
  spec:
    infrastructureRef:
      apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
      kind: DockerCluster
      name: #{name}
  YAML

  IO.popen('kubectl apply -f -', 'w') do |io|
    io.puts yaml
  end
  puts "Created cluster #{name}"
end

if ARGV[0] == 'list'
  list_clusters
elsif ARGV[0] == 'create'
  create_demo_cluster(ARGV[1] || 'demo-cluster')
else
  puts "Usage: #{__FILE__} [list|create <name>]"
end
```

---

## File: `docs/quickstart.md`

```markdown
# Quickstart — dev flow

1. On your Mac or dev machine, ensure Docker, kubectl, clusterctl are installed.
2. `./bootstrap/00-check-prereqs.sh`
3. `./bootstrap/01-install-clusterctl.sh`
4. Create a dev management cluster (kind or docker) and `clusterctl init --infrastructure docker`
5. Build backend and run locally: `make build && make docker-build`
6. Run the backend container with `make docker-run` and use Postman or `curl` to call:
   - `GET http://localhost:8080/healthz`
   - `POST http://localhost:8080/v1/clusters` with `{"name":"demo1","namespace":"default"}`

Troubleshooting: If your kubeconfig is in a nonstandard path, set `KUBECONFIG` accordingly before running.
```

---

## Final notes
- This repository scaffold intentionally keeps things small so you can iterate quickly. For production: add unit tests, integration tests, proper logging, structured errors, observability, tracing, authentication, and RBAC.
- If you want, I can convert this into individual downloadable files (zip) or generate additional languages (Python, Rust) for specific components.

---

*End of scaffold.*

