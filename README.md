# Kubernetes Microservices Deployment

Declarative Kubernetes manifests for deploying an independent, scalable 4-tier microservices architecture.

---

## Architecture Overview

```text
apps/                       # Go microservices source code & multi-stage Dockerfiles
├── user-service/           # User authentication & profiles (~15MB image)
├── inventory-service/      # Stock & availability
├── media-service/          # Video streaming & assets
└── lesson-service/         # Microservice orchestrator (aggregates data via cluster DNS)

services/                   # Declarative Kubernetes manifests (Deployments & Services)
├── user-service.yaml       # User authentication & profiles (3 replicas)
├── inventory-service.yaml  # Stock & availability (2 replicas)
├── media-service.yaml      # Video streaming & assets (3 replicas)
└── lesson-service.yaml     # Course curriculum & progress (2 replicas)
```


Each service is decoupled with its own:
- **Deployment**: Horizontal replicas, container specifications, resource limits & requests.
- **Service**: Dedicated `ClusterIP` front-door with internal CoreDNS discovery.

---

## Deploying to Kubernetes

### 1. Deploy all microservices in parallel
```bash
kubectl apply -f services/
```

### 2. Verify deployments and services
```bash
kubectl get deployments,pods,services -l tier=backend
```

### 3. Test inter-service communication
```bash
kubectl run test-curl --rm -it --image=curlimages/curl -- curl http://user-service
```
